# docker-credential-vault

A [Docker credential helper](https://github.com/docker/docker-credential-helpers) to store credentials in [HashiCorp Vault](https://vaultproject.io).

## Usage

 1. In your docker config (`~/.docker/config.json` typically), set:

    ```json
    {
        "credsStore": "vault"
    }
    ```

 1. Ensure `docker-credential-vault` is in your `$PATH`.
 1. Decide which Vault KV backend (both v1 and v2 are supported) and the path t
 1. `export DOCKER_CREDENTIAL_VAULT_KV_PATH=secret/path/to/use`: This is a KV backend (both v1 and v2 are supported) path where the helper will store and look for credentials.
 1. Configure how the helper will connect to vault, all the standard `$VAULT_` environment variables will be used to configured the vault client, e.g. `$VAULT_ADDR` and `$VAULT_TOKEN` will be required at a minimum. 

## Internals

The helper will store the credentials in the following format, e.g. server URL is `https://example.com:8080`:

`vault kv get -format=json $DOCKER_CREDENTIAL_VAULT_KV_PATH/example.com` will return:

```json
{
  "request_id": "89ec7fd0-be41-2e85-7c79-1b16199a3d7b",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "data": {
      "Secret": "s3cr3t",
      "ServerURL": "https://example.com",
      "Username": "docker"
    },
    "metadata": {"redacted": "metadata"}
  },
  "warnings": null
}
```
