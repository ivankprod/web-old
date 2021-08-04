echo "Starting server building..."

env_contents=$(<../../.env)

mode="prod"
if [[ "$1" == "dev" ]]; then
	mode=$1

	echo "${env_contents//STAGE_MODE=\"prod\"/STAGE_MODE=\"dev\"}" > ../../.env
elif [[ "$1" == "prod" ]] || [[ "$1" == "" ]]; then
	mode="prod"

	echo "${env_contents//STAGE_MODE=\"dev\"/STAGE_MODE=\"prod\"}" > ../../.env
else
	echo "Error argument parsing"; exit
fi

echo "MODE: $mode"

os="$2"
arch="$3"

windres -o server-res.syso ./resources_win/server.rc
pkger
GOOS=$os GOARCH=$arch GOMAXPROCS=4 go build -o ../../build_$mode -v -ldflags="-s -w"

cp ../../.env ../../build_$mode/
mkdir -p ../../build_$mode/logs

echo "Backend build done!"