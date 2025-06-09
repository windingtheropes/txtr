build:
	GOOS=linux
	GOARCH=amd64 go build -o dist/linux-amd64
	GOARCH=arm64 go build -o dist/linux-arm64
	GOOS=darwin
	GOARCH=amd64 go build -o dist/mac-amd64
	GOARCH=arm64 go build -o dist/mac-arm64
	GOOS=windows
	GOARCH=amd64 go build -o dist/windows-amd64.exe
	GOARCH=arm64 go build -o dist/windows-arm64.exe