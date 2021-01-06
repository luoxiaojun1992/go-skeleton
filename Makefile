DARWIN_DIR = bin/darwin_amd64
LINUX_DIR = bin/linux_amd64
WINDOWS_DIR = bin/windows_amd64

darwin:
	GOOS=darwin GOARCH=amd64 go build -i -v -o ./${DARWIN_DIR}/command

linux:
	GOOS=linux GOARCH=amd64 go build -i -v -o ./${LINUX_DIR}/command

windows:
	GOOS=windows GOARCH=amd64 go build -i -v -o ./${WINDOWS_DIR}/command