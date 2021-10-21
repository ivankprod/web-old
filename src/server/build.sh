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

if [[ $os == "windows" ]]; then
	windres -o server-res.syso ./resources_win/server.rc
fi

pkger
GOOS=$os GOARCH=$arch GOMAXPROCS=4 go build -o ../../build_$mode -v -ldflags="-s -w"

cp ../../$mode.env ../../build_$mode/.env
mkdir -p ../../build_$mode/logs

echo "Backend build done!"