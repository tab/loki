# Certificates and keys

## Generating JWT Signing Key pair

Loki uses JWT (JSON Web Tokens) for authentication.

To generate a signing key for JWT, you can use the following command:

```sh
openssl genrsa -out certs/jwt/private.key 4096
```

To generate the public key from the private key, use:

```sh
openssl rsa -in certs/jwt/private.key -pubout -out certs/jwt/public.key
```

## Generating Certificates for mTLS

For mTLS (mutual TLS), both the server and client need certificates.
The process involves:

- Creating a Certificate Authority (CA)
- Creating server certificates signed by the CA
- Creating client certificates signed by the CA

### Generate the Certificate Authority (CA)

Generate a private key for your CA

```sh
openssl genrsa -out certs/ca.key 4096
openssl req -new -x509 -key certs/ca.key -sha256 -subj '/CN=Loki CA' -out certs/ca.pem -days 3650
```

### Generate the Server Certificate

#### Generate server private key

```sh
openssl genrsa -out certs/server.key 4096
```

#### Create server Certificate Signing Request (CSR)

```sh
openssl req -new -key certs/server.key -out certs/server.csr -config <(
cat <<-EOF
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
IP.1 = 127.0.0.1
IP.2 = 0.0.0.0
EOF
)
```

#### Sign the server certificate with CA

```sh
openssl x509 -req -in certs/server.csr -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -out certs/server.pem -days 825 -sha256 -extfile <(
cat <<-EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = backend
IP.1 = 127.0.0.1
IP.2 = 0.0.0.0
EOF
)
```

### Generate the Client Certificate

Generate client private key

```sh
openssl genrsa -out certs/client.key 4096
```

#### Create client Certificate Signing Request (CSR)

```sh
openssl req -new -key certs/client.key -out certs/client.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn

[dn]
CN = loki-backoffice
EOF
)
```

#### Sign the client certificate with CA

```sh
openssl x509 -req -in certs/client.csr -CA certs/ca.pem -CAkey certs/ca.key -CAcreateserial -out certs/client.pem -days 825 -sha256
```

### Verify the certificates

```sh
openssl verify -CAfile certs/ca.pem certs/server.pem
```

```sh
openssl verify -CAfile certs/ca.pem certs/client.pem
```
