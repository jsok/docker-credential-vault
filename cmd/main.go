package main

import (
	"fmt"
	"os"

	"github.com/docker/docker-credential-helpers/credentials"
	"github.com/hashicorp/vault/api"
	helper "github.com/jsok/docker-credential-vault"
)

const kvPathEnvVar = "DOCKER_CREDENTIAL_VAULT_KV_PATH"

func main() {
	kvPath := os.Getenv(kvPathEnvVar)
	if kvPath == "" {
		fmt.Fprintf(os.Stdout, "%s must be defined\n", kvPathEnvVar)
		os.Exit(1)
	}
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		fmt.Fprintf(os.Stdout, "Could not create vault client: %v\n", err)
		os.Exit(1)
	}
	credentials.Serve(helper.NewHelper(kvPath, client))
}
