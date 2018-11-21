export GOARCH=amd64
export GOOS=linux

export CURR=`pwd`

cd ../../../../../..

export GOPATH=`pwd`

cd ${CURR}

echo "GOPATH : $GOPATH"
echo "GOROOT : $GOROOT"
echo "GOARCH : $GOARCH"
echo "GOOS : $GOOS"

go build -o ../linux/messageProtoJson github.com/Blizzardx/MessageProtoJson/release
