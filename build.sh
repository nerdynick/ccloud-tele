#Windows 64bit
export GOOS=windows
export GOARCH=amd64
go build -o ccloud-tele-win-64

#Linux 64bit
export GOOS=linux
export GOARCH=amd64
go build -o ccloud-tele-linux-64

#Mac 64bit
export GOOS=darwin
export GOARCH=amd64
go build -o ccloud-tele-mac-64

#Mac ARM 64bit
export GOOS=darwin
export GOARCH=arm64
go build -o ccloud-tele-mac-arm64