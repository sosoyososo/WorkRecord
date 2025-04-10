env GOOS=linux GOARCH=amd64 go build -o linux_amd64 main.go;
env GOOS=darwin GOARCH=amd64 go build -o mac_amd64 main.go;
env GOOS=darwin GOARCH=arm64 go build -o mac_arm64 main.go;
