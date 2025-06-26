#!/bin/bash

set -e

echo "🛠️  Starting PostgreSQL scaffold setup..."

# Prompt for admin password
read -sp "🔐 Enter a password for the admin PostgreSQL user: " ADMIN_PASSWORD
echo ""

# Generate secure API user password
API_PASSWORD=$(openssl rand -base64 12)
echo "✅ API user password generated."

# Create .env.db
echo "📁 Creating .env.db..."
cat <<EOF >.env.db
POSTGRES_DB=saas_api_dev
POSTGRES_USER=postgres
POSTGRES_PASSWORD=$ADMIN_PASSWORD
API_USER=api_user
API_PASSWORD=$API_PASSWORD
EOF
echo "✅ .env.db created."

# Create docker-compose.yml
echo "📁 Creating docker-compose.yml..."
cat <<EOF >docker-compose.yml
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
echo "✅ docker-compose.yml created."

# Create init SQL directory and file
echo "📁 Creating init SQL script..."
mkdir -p init
source ../environments/.env.db

cat <<EOF >./postgres/create-api-user.sql
DO
\$\$
BEGIN
   IF NOT EXISTS (
      SELECT FROM pg_catalog.pg_roles
      WHERE rolname = '${API_USER}') THEN
      CREATE ROLE ${API_USER} LOGIN PASSWORD '${API_PASSWORD}';
   END IF;
END
\$\$;

GRANT CONNECT ON DATABASE ${POSTGRES_DB} TO ${API_USER};
GRANT USAGE ON SCHEMA public TO ${API_USER};
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO ${API_USER};
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO ${API_USER};
EOF
echo "✅ sql/create-api-user.sql created."

# Final prompt and optional DB reset
echo ""
echo "🎯 Scaffold complete."
read -p "🔄 Reset existing DB container and volume now? (y/N): " RESET

if [[ "$RESET" =~ ^[Yy]$ ]]; then
    echo "🧹 Cleaning up existing container and volume..."
    docker compose down -v
    echo "🚀 Starting fresh container with updated credentials..."
    docker compose --env-file ../environments/.env.db up -d
    echo "✅ Database reset and container restarted."
else
    echo "ℹ️  You can start the container manually with:"
    echo "    docker compose --env-file ./environments/.env.db up -d"
fi
