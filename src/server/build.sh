echo "Starting server building..."

mode="prod"
if [[ "$1" == "dev" ]]; then
	mode=$1
elif [[ "$1" == "prod" ]] || [[ "$1" == "" ]]; then
	mode="prod"
else
	echo "Error argument parsing"; exit
fi

echo "MODE: $mode"

os="$2"
arch="$3"

GOOS=$os GOARCH=$arch GOMAXPROCS=4 go build -o ../../build_$mode/server -v -ldflags="-s -w" ./cmd/main.go

mkdir -p ../../build_$mode/misc
cp ./misc/sitemap.json ../../build_$mode/misc/sitemap.json

cp ../../$mode.env ../../build_$mode/.env

mkdir -p ../../build_$mode/logs

echo "Server build done!"