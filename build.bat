go get github.com/akavel/rsrc
rsrc -arch 386 -ico icon.ico -o rsrc_386.syso
rsrc -arch amd64 -ico icon.ico -o rsrc_amd64.syso

set GOARCH=386
go build -ldflags "-s -w -H=windowsgui" -o jolina.exe