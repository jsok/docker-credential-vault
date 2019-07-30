# docker-credential-vault

A [Docker credential helper](https://github.com/docker/docker-credential-helpers) to store credentials in [HashiCorp Vault](https://vaultproject.io).

## KV setup

The helper will require a KV backend configured in the Vault server.

In order to configure the helper you will need to:

 1. `export DOCKER_CREDENTIAL_VAULT_KV_PATH=secret/path/to/use`: This is a KV backend path where the helper will store and look for credentials.
 1.  All the standard `$VAULT_` environment variables will be used to configured the vault client, e.g. `VAULT_ADDR`, `VAULT_TOKEN` etc.

## Internal

The helper will store the credentials in the following format, e.g. server URL is `https://example.com:8080`:

`vault read -format=json $DOCKER_CREDENTIAL_VAULT_KV_PATH/example.com`

That secret will contain the data body:

```json
{
  "ServerURL": "https://example.com:8080",
  "Username": "username",
  "Secret": "secret"
}
```
