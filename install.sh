echo "Installing webapp..."

echo "Installing frontend"
cd src/frontend/
npm install

echo "Installing backend"
cd ../server/
go get

cd ../../
cp .env.local .env

echo "Installing webapp done!"