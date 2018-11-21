set GOARCH=amd64
set GOOS=windows
set CURR=%cd%
cd ../../../../../../

set GOPATH=%cd%
cd %CURR%

go build -o ../windows/messageProtoJson.exe github.com/Blizzardx/MessageProtoJson/release

@IF %ERRORLEVEL% NEQ 0 pause

