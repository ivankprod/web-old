echo "Starting frontend building..."

mode="prod"
env="production"

if [[ "$1" == "dev" ]]; then
	mode=$1
	env="development"
elif [[ "$1" == "prod" ]] || [[ "$1" == "" ]]; then
	mode="prod"
	env="production"
else
	echo "Error argument parsing"; exit
fi

echo "MODE: $mode"

cp -r ../server/views/ ../../build_$mode/

NODE_ENV="$env" npm run "$mode"

echo "Frontend build done!"