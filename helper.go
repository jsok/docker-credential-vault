package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/docker/docker-credential-helpers/registryurl"
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
	client := h.client.Logical()
	path, err := h.pathForServerURL(creds.ServerURL)
	if err != nil {
		return err
	}
	_, err = client.Write(path, map[string]interface{}{
		"ServerURL": creds.ServerURL,
		"Username":  creds.Username,
		"Secret":    creds.Secret,
	})
	return err
}

func (h *Helper) Delete(serverURL string) error {
	client := h.client.Logical()
	path, err := h.pathForServerURL(serverURL)
	if err != nil {
		return err
	}
	_, err = client.Delete(path)
	return err
}

func (h *Helper) Get(serverURL string) (string, string, error) {
	creds, err := h.read(serverURL)
	if err != nil {
		return "", "", err
	}
	return creds.Username, creds.Secret, nil
}

func (h *Helper) List() (map[string]string, error) {
	client := h.client.Logical()

	secret, err := client.List(h.kvPath)
	if err != nil {
		return nil, err
	}

	keys, ok := extractListData(secret)
	if !ok {
		return nil, fmt.Errorf("No credentials found at %s\n", h.kvPath)
	}

	var result map[string]string = make(map[string]string, 0)
	for _, key := range keys {
		server := key.(string)
		if strings.HasSuffix(server, "/") {
			// Skip any sub-paths, only looking for leaf keys
			continue
		}
		creds, err := h.read(server)
		if err != nil {
			continue
		}
		result[creds.ServerURL] = creds.Username
	}

	return result, nil
}

// Internal helper to read a credential from vault
func (h *Helper) read(serverURL string) (*credentials.Credentials, error) {
	client := h.client.Logical()
	path, err := h.pathForServerURL(serverURL)
	if err != nil {
		return nil, err
	}
	secret, err := client.Read(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", path, err)
	}

	return secretToCredential(secret)
}

func (h *Helper) pathForServerURL(serverURL string) (string, error) {
	u, err := registryurl.Parse(serverURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", h.kvPath, registryurl.GetHostname(u)), nil
}

// Convert a Vault Secret to a Docker Credential
func secretToCredential(secret *api.Secret) (*credentials.Credentials, error) {
	server, exists := secret.Data["ServerURL"]
	if !exists {
		return nil, credentials.NewErrCredentialsMissingServerURL()
	}
	username, exists := secret.Data["Username"]
	if !exists {
		return nil, credentials.NewErrCredentialsMissingUsername()
	}
	password, exists := secret.Data["Secret"]
	if !exists {
		return nil, fmt.Errorf("no credentials secret")
	}
	return &credentials.Credentials{
		ServerURL: server.(string),
		Username:  username.(string),
		Secret:    password.(string),
	}, nil
}

// extractListData reads the secret and returns a typed list of data and a
// boolean indicating whether the extraction was successful.
func extractListData(secret *api.Secret) ([]interface{}, bool) {
	if secret == nil || secret.Data == nil {
		return nil, false
	}

	k, ok := secret.Data["keys"]
	if !ok || k == nil {
		return nil, false
	}

	i, ok := k.([]interface{})
	return i, ok
}
