export POSTGRES_CONN_STRING=postgresql://postgres:postgres@localhost/demo?sslmode=disable
export REACT_APP_ENVIRONMENT="development"
export NODE_ENV="development"

pushd _ui
npm run build:dev
popd
go run .