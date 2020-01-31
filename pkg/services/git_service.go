package services

import (
	"errors"
	"fmt"
	"github.com/saas/hostgolang/pkg/config"
	"log"
	"net/http"
	"time"
)

type GitService interface {
	CreateRepository(name string) error
	WriteNotification(key, message string)
}

type httpGitService struct {
	httpClient *http.Client
	cfg *config.Config
	tcp *TcpServer
}

func NewGitService(cfg *config.Config) GitService {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	tcp := NewTcpServer()
	go func() {
		log.Println("git TCP notification server started at ", cfg.GitTcpAddr)
		if err := tcp.Run(cfg.GitTcpAddr); err != nil {
			log.Fatal("failed to start Git tcp server: ", err)
		}
	}()
	return &httpGitService{httpClient: httpClient, cfg: cfg, tcp: tcp}
}

func (h *httpGitService) WriteNotification(key, message string) {
	data := []byte(message)
	if err := h.tcp.Send(key, data); err != nil {
		log.Println("failed to write TCP message: ", err)
	}
}

func (h *httpGitService) CreateRepository(name string) error {
	u := fmt.Sprintf("%s/repo/new?name=%s", h.cfg.GitServerUrl, name)
	r, err := http.Post(u, "application/json", nil)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New("failed to create repository at this time")
	}
	return nil
}
