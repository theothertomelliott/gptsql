export POSTGRES_CONN_STRING=postgresql://postgres:postgres@localhost/demo?sslmode=disable
export USE_DEV_FRONTEND="true"

pushd _ui
npm run build
popd
go run .