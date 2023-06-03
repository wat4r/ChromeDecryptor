set VERSION=1.0.0

@REM Linux
@REM set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -trimpath -o bin/chromeDecryptor_%VERSION%_%GOOS%_%GOARCH% main.go

@REM Mac
@REM SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -ldflags "-s -w" -trimpath -o bin/chromeDecryptor_%VERSION%_%GOOS%_%GOARCH% main.go

@REM Windows
@REM set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -trimpath -o bin/chromeDecryptor_%VERSION%_%GOOS%_%GOARCH%.exe main.go
