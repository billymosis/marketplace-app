#! /bin/bash
migrate -path ./migrations -database "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable" up
