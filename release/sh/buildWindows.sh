export GOARCH=amd64
export GOOS=windows

export CURR=`pwd`
cd ../../../../../..

export GOPATH=`pwd`

cd ${CURR}

echo "GOPATH : $GOPATH"
echo "GOROOT : $GOROOT"
echo "GOARCH : $GOARCH"
echo "GOOS : $GOOS"

go build -o ../windows/messageProtoJson.exe github.com/Blizzardx/MessageProtoJson/release
