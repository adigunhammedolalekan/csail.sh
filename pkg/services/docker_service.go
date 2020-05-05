package services

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/google/uuid"
	"github.com/saas/hostgolang/pkg/config"
	"io"
	"log"
	"os"
	"path/filepath"
)

const baseBuildDir = "/tmp/mnt/hostgo/apps"
const rawDockerfile = `
FROM alpine:3.2
RUN apk update && apk add --no-cache ca-certificates
ADD . /app
WORKDIR /app
RUN chmod +x /app/%s
ENTRYPOINT [ "/app/%s" ]`

//go:generate mockgen -destination=../mocks/docker_service_mock.go -package=mocks github.com/saas/hostgolang/pkg/services DockerService
type DockerService interface {
	CopyToContainer(destFileName, containerId string, content io.Reader) error
}

type defaultDockerService struct {
	client *client.Client
	cfg    *config.Config
}

func NewDockerService(cli *client.Client, cfg *config.Config) DockerService {
	return &defaultDockerService{client: cli, cfg: cfg}
}

func (d *defaultDockerService) CopyToContainer(destFileName, containerId string, content io.Reader) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	f := filepath.Join(wd, destFileName)
	fi, err := os.Create(f)
	if err != nil {
		return err
	}
	if _, err := io.Copy(fi, content); err != nil {
		return err
	}
	tarredContent, err := d.createBuildContext(f)
	if err != nil {
		return err
	}
	dstDir := filepath.Join("mnt", "tmp")
	err = d.client.CopyToContainer(context.Background(),
		containerId, dstDir, tarredContent, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
	if err != nil {
		log.Println("cannot copy to container ", err)
		return err
	}
	// clean up
	return os.Remove(f)
}

func (d *defaultDockerService) createBuildContext(filename string) (io.Reader, error) {
	return archive.Tar(filename, archive.Uncompressed)
}

func (d *defaultDockerService) md5() string {
	m5 := md5.New()
	m5.Write([]byte(uuid.New().String()))
	return fmt.Sprintf("%+x", string(m5.Sum(nil)))
}

func (d *defaultDockerService) registryAuthAsBase64() string {
	authConfig := types.AuthConfig{
		Username: d.cfg.Registry.Username,
		Password: d.cfg.Registry.Password,
	}
	encoded, err := json.Marshal(authConfig)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(encoded)
}
