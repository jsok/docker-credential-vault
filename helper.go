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
	data := map[string]interface{}{
		"ServerURL": creds.ServerURL,
		"Username":  creds.Username,
		"Secret":    creds.Secret,
	}
	path, err := h.pathForServerURL(creds.ServerURL)
	if err != nil {
		return err
	}

	mountPath, v2, err := isKVv2(path, h.client)
	if err != nil {
		return err
	}
	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		data = map[string]interface{}{
			"data":    data,
			"options": map[string]interface{}{},
		}
	}
	_, err = h.client.Logical().Write(path, data)
	return err
}

func (h *Helper) Delete(serverURL string) error {
	path, err := h.pathForServerURL(serverURL)
	if err != nil {
		return err
	}
	mountPath, v2, err := isKVv2(path, h.client)
	if err != nil {
		return err
	}
	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
	}
	_, err = h.client.Logical().Delete(path)
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
	path := h.kvPath
	mountPath, v2, err := isKVv2(path, h.client)
	if err != nil {
		return nil, err
	}
	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "metadata")
	}

	secret, err := h.client.Logical().List(path)
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
// which understands both KV v1 and v2 APIs.
func (h *Helper) read(serverURL string) (*credentials.Credentials, error) {
	path, err := h.pathForServerURL(serverURL)
	if err != nil {
		return nil, err
	}

	mountPath, v2, err := isKVv2(path, h.client)
	if err != nil {
		return nil, err
	}
	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
	}
	secret, err := kvReadRequest(h.client, path, nil)
	if err != nil || secret == nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", path, err)
	}

	data := secret.Data
	if v2 && data != nil {
		data = nil
		dataRaw := secret.Data["data"]
		if dataRaw != nil {
			data = dataRaw.(map[string]interface{})
		}
	}
	return secretDataToCredential(data)
}

func (h *Helper) pathForServerURL(serverURL string) (string, error) {
	u, err := registryurl.Parse(serverURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", h.kvPath, registryurl.GetHostname(u)), nil
}

// Convert a Vault Secret to a Docker Credential
func secretDataToCredential(data map[string]interface{}) (*credentials.Credentials, error) {
	server, exists := data["ServerURL"]
	if !exists {
		return nil, credentials.NewErrCredentialsMissingServerURL()
	}
	username, exists := data["Username"]
	if !exists {
		return nil, credentials.NewErrCredentialsMissingUsername()
	}
	password, exists := data["Secret"]
	if !exists {
		return nil, fmt.Errorf("no credentials secret")
	}
	return &credentials.Credentials{
		ServerURL: server.(string),
		Username:  username.(string),
		Secret:    password.(string),
	}, nil
}
