#! /bin/bash

cleanup() {
    echo "Performing cleanup..."
    docker compose down
    echo "Cleanup complete."
}

trap cleanup EXIT

docker compose up -d
air
