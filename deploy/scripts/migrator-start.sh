#!/bin/sh
export MIGRATION_DSN="host=${PG_HOST} port=${PG_PORT} dbname=${POSTGRES_DB} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} sslmode=disable"

goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v