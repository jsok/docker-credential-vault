package helper

import (
	"github.com/docker/docker-credential-helpers/credentials"
)

type Helper struct{}

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
