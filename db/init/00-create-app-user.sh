#!/usr/bin/env bash
set -e

echo "==> Creating application user in admin..."
mongosh --quiet --username "${MONGO_INITDB_ROOT_USERNAME}" \
        --password "${MONGO_INITDB_ROOT_PASSWORD}" \
        --authenticationDatabase "admin" <<'EOF'
use admin
db.createUser({
  user: process.env.MONGO_USERNAME || "app_user",
  pwd:  process.env.MONGO_PASSWORD || "app_pass_123",
  roles: [
    { role: "readWrite", db: process.env.MONGO_INITDB_DATABASE || "local" }
  ]
});
EOF
echo "==> App user created ✅"
