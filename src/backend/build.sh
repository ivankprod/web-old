echo "Starting backend building..."

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

windres -o server-res.syso ./resources_win/server.rc
pkger
GOOS=windows GOARCH=amd64 GOMAXPROCS=4 go build -o ../../build_$mode/server.exe -v -ldflags="-s -w"

echo "Backend build done!"