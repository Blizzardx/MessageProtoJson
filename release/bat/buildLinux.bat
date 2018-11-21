set GOARCH=amd64
set GOOS=linux
set CURR=%cd%
cd ../../../../../../

set GOPATH=%cd%
cd %CURR%

go build -o ../linux/messageProtoJson github.com/Blizzardx/MessageProtoJson/release

@IF %ERRORLEVEL% NEQ 0 pause

