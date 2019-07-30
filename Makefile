bin:
	mkdir -p $@

vault: | bin
	go build -ldflags -s -o bin/docker-credential-vault ./cmd/main.go
