echo "Starting building webapp..."

mode="prod"
if [[ "$1" == "dev" ]]; then
	mode=$1
elif [[ "$1" == "prod" ]] || [[ "$1" == "" ]]; then
	mode="prod"
else
	echo "Error argument parsing"; exit
fi

cd src/frontend/
./build.sh "$mode"

cd ../server/
./build.sh "$mode"

cd ../../

mkdir -p build_"$mode"/logs

echo "Building webapp done!"