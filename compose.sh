echo "Starting composing webapp..."

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

. ./build_$mode/.env

set -a
STAGE_MODE=$STAGE_MODE
NODE_ENV=$env
SERVER_HOST=$SERVER_HOST
SERVER_PORT_HTTP=$SERVER_PORT_HTTP
SERVER_PORT_HTTPS=$SERVER_PORT_HTTPS
DB_TARANTOOL_PORT=$DB_TARANTOOL_PORT
DB_TARANTOOL_USER=$DB_TARANTOOL_USER
DB_TARANTOOL_PASSWORD=$DB_TARANTOOL_PASSWORD
GRAFANA_ADMIN_USER=$GRAFANA_ADMIN_USER
GRAFANA_ADMIN_PASSWORD=$GRAFANA_ADMIN_PASSWORD
set +a

if [[ "$2" != "--nobuild" ]]; then
	docker-compose build app
fi

docker-compose up -d

echo "Composing webapp done!"