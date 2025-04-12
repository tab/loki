#!/bin/bash
set -e

# This script generates JWT keys and mTLS certificates for Loki and Loki-backoffice
# Usage: ./generate-certs.sh [LOKI_REPO] [LOKI_BACKOFFICE_REPO]

LOKI_REPO="${1:-../../loki}"
LOKI_BACKOFFICE_REPO="${2:-../../loki-backoffice}"

echo "Loki repo: $LOKI_REPO"
echo "Loki-backoffice repo: $LOKI_BACKOFFICE_REPO"

mkdir -p "$LOKI_REPO/certs/jwt"
mkdir -p "$LOKI_BACKOFFICE_REPO/certs/jwt"

echo "Generating JWT signing keys..."
openssl genrsa -out "$LOKI_REPO/certs/jwt/private.key" 4096
openssl rsa -in "$LOKI_REPO/certs/jwt/private.key" -pubout -out "$LOKI_REPO/certs/jwt/public.key"
cp "$LOKI_REPO/certs/jwt/public.key" "$LOKI_BACKOFFICE_REPO/certs/jwt/public.key"

echo "Generating Certificate Authority (CA)..."
openssl genrsa -out "$LOKI_REPO/certs/ca.key" 4096
openssl req -new -x509 -key "$LOKI_REPO/certs/ca.key" -sha256 -subj "/CN=Loki CA" \
    -out "$LOKI_REPO/certs/ca.pem" -days 3650
cp "$LOKI_REPO/certs/ca.pem" "$LOKI_BACKOFFICE_REPO/certs/ca.pem"

echo "Creating temporary config files..."
cat > server_config.cnf << EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
CN = loki-backend

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = backend
DNS.3 = loki
DNS.4 = loki-backend
IP.1 = 127.0.0.1
IP.2 = 0.0.0.0
EOF

cat > server_ext.cnf << EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = backend
DNS.3 = loki
DNS.4 = loki-backend
IP.1 = 127.0.0.1
IP.2 = 0.0.0.0
EOF

cat > client_config.cnf << EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn

[dn]
CN = loki-backoffice
EOF

echo "Generating Server Certificate..."
openssl genrsa -out "$LOKI_REPO/certs/server.key" 4096
openssl req -new -key "$LOKI_REPO/certs/server.key" \
    -out "$LOKI_REPO/certs/server.csr" -config server_config.cnf
openssl x509 -req -in "$LOKI_REPO/certs/server.csr" -CA "$LOKI_REPO/certs/ca.pem" \
    -CAkey "$LOKI_REPO/certs/ca.key" -CAcreateserial \
    -out "$LOKI_REPO/certs/server.pem" -days 825 -sha256 -extfile server_ext.cnf

echo "Generating Client Certificate..."
openssl genrsa -out "$LOKI_BACKOFFICE_REPO/certs/client.key" 4096
openssl req -new -key "$LOKI_BACKOFFICE_REPO/certs/client.key" \
    -out "$LOKI_BACKOFFICE_REPO/certs/client.csr" -config client_config.cnf
openssl x509 -req -in "$LOKI_BACKOFFICE_REPO/certs/client.csr" -CA "$LOKI_REPO/certs/ca.pem" \
    -CAkey "$LOKI_REPO/certs/ca.key" -CAcreateserial \
    -out "$LOKI_BACKOFFICE_REPO/certs/client.pem" -days 825 -sha256

echo "Verifying certificates..."
openssl verify -CAfile "$LOKI_REPO/certs/ca.pem" "$LOKI_REPO/certs/server.pem"
openssl verify -CAfile "$LOKI_REPO/certs/ca.pem" "$LOKI_BACKOFFICE_REPO/certs/client.pem"

rm -f server_config.cnf server_ext.cnf client_config.cnf

echo "Generated JWT keys and certificates:"
echo "Loki certificates:"
ls -la "$LOKI_REPO/certs"
ls -la "$LOKI_REPO/certs/jwt"
echo "Loki-backoffice certificates:"
ls -la "$LOKI_BACKOFFICE_REPO/certs"
ls -la "$LOKI_BACKOFFICE_REPO/certs/jwt"
