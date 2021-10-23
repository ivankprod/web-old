echo "Starting composing webapp..."

mode="prod"
if [[ "$1" == "dev" ]]; then
	mode=$1
elif [[ "$1" == "prod" ]] || [[ "$1" == "" ]]; then
	mode="prod"
else
	echo "Error argument parsing"; exit
fi

echo "MODE: $mode"

. ./$mode.env

set -a
STAGE_MODE=$STAGE_MODE
DB_TARANTOOL_PORT=$DB_TARANTOOL_PORT
DB_TARANTOOL_USER=$DB_TARANTOOL_USER
DB_TARANTOOL_PASSWORD=$DB_TARANTOOL_PASSWORD
set +a

docker-compose up -d

echo "Composing webapp done!"