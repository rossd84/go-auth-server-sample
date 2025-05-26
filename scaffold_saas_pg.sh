#!/bin/bash

# Exit on any error
set -e

# Ask for admin password
read -sp "Enter a password for the admin PostgreSQL user: " ADMIN_PASSWORD
echo

# Generate secure API user password
API_PASSWORD=$(openssl rand -base64 12)
echo "Generated API password: $API_PASSWORD"

# Create environment file
cat <<EOF >.env.db
POSTGRES_DB=saas_api_dev
POSTGRES_USER=postgres
POSTGRES_PASSWORD=$ADMIN_PASSWORD
API_USER=api_user
API_PASSWORD=$API_PASSWORD
EOF

# Create docker-compose.yml
cat <<EOF >docker-compose.yml
version: '3.8'

services:
  db:
    image: postgres:latest
    container_name: saas_postgres_dev
    restart: unless-stopped
    env_file:
      - .env.db
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./init:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"

volumes:
  pgdata:
EOF

# Create init directory if not exists
mkdir -p init

# Create SQL script to add API user
cat <<EOF >init/create-api-user.sql
DO
\\$\\$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = 'api_user') THEN
      CREATE ROLE api_user LOGIN PASSWORD '\${API_PASSWORD}';
   END IF;
END
\\$\\$;

GRANT CONNECT ON DATABASE saas_api_dev TO api_user;
GRANT USAGE ON SCHEMA public TO api_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO api_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO api_user;
EOF

echo "✅ Scaffold complete. You can now run: docker compose up -d"
