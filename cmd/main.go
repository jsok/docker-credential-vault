package main

import (
	"github.com/docker/docker-credential-helpers/credentials"
	vaultHelper "github.com/jsok/docker-credential-vault"
)

func main() {
	credentials.Serve(&vaultHelper.Helper{})
}
