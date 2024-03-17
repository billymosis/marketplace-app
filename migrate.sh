#! /bin/bash
migrate -path ./migrations/db -database "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable" up
