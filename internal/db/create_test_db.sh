#!/bin/bash

# Define the database file
DB_FILE="/tmp/spacetraders.db"

# Define the schema as a variable
USERSCHEMA=$(cat <<EOF
CREATE TABLE users (
    username TEXT PRIMARY KEY,
    password BLOB,
    faction TEXT,
    accessToken TEXT
);
EOF
)

REFRESHSCHEMA=$(cat <<EOF
CREATE TABLE refreshTokens (
    refreshToken TEXT PRIMARY KEY,
    accessToken TEXT,
    createdAt INTEGER
);
EOF
)

# Execute the schema to create the database and table
echo "$USERSCHEMA" | sqlite3 "$DB_FILE"
echo "$REFRESHSCHEMA" | sqlite3 "$DB_FILE"
echo "Database and tables created successfully."
