docker run --name postgresdb -v $(pwd)/init/database:/docker-entrypoint-initdb.d -p 5432:5432 -e POSTGRES_USER=testapi -e POSTGRES_PASSWORD=fizznbuzz -e POSTGRES_DB=testdb -d postgres:14.0
