#!/usr/bin/env bash

set -euo pipefail

DB_NAME="${DB_NAME:-animalsys}"
MONGO_SERVICE="${MONGO_SERVICE:-mongodb}"
BACKEND_SERVICE="${BACKEND_SERVICE:-backend}"

echo "🗑️  Dropping MongoDB database '${DB_NAME}' (service: ${MONGO_SERVICE})..."
docker-compose exec "${MONGO_SERVICE}" mongosh "${DB_NAME}" --eval "db.dropDatabase()"

echo "🌱 Running seed script from ${BACKEND_SERVICE}..."
docker-compose exec "${BACKEND_SERVICE}" go run ./cmd/seed

echo "✅ Database reseeded successfully."
