package helper

import (
	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/hashicorp/vault/api"
)

type Helper struct {
	kvPath string
	client *api.Client
}

func NewHelper(kvPath string, client *api.Client) *Helper {
	return &Helper{kvPath, client}
}

func (h *Helper) Add(creds *credentials.Credentials) error {
	return nil
}

func (h *Helper) Delete(serverURL string) error {
	return nil
}

func (h *Helper) Get(serverURL string) (string, string, error) {
	return "", "", nil
}

func (h *Helper) List() (map[string]string, error) {
	return nil, nil
}
