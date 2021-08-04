echo "Starting building webapp..."

mode="prod"
if [[ "$1" == "dev" ]]; then
	mode=$1
elif [[ "$1" == "prod" ]] || [[ "$1" == "" ]]; then
	mode="prod"
else
	echo "Error argument parsing"; exit
fi

os="$2"
arch="$3"

cd src/frontend/
./build.sh "$mode"

cd ../server/
./build.sh "$mode" "$os" "$arch"

cd ../../

echo "Building webapp done!"