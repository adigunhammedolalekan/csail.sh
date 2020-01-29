package services

import (
	"errors"
	"fmt"
	"github.com/saas/hostgolang/pkg/config"
	"net/http"
	"time"
)

type GitService interface {
	CreateRepository(name string) error
}

type httpGitService struct {
	httpClient *http.Client
	cfg *config.Config
}

func NewGitService(cfg *config.Config) GitService {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	return &httpGitService{httpClient: httpClient, cfg: cfg}
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
