#!/bin/bash
# this was generated by chatGPT

# Generate private key
openssl genrsa -out https-example/localhost.key 2048

# Generate Certificate Signing Request (CSR)
openssl req -new -key https-example/localhost.key -out https-example/localhost.csr -subj "/CN=localhost"

# Generate self-signed certificate valid for 365 days
openssl x509 -req -days 365 -in https-example/localhost.csr -signkey https-example/localhost.key -out https-example/localhost.crt

# Cleanup CSR file
rm https-example/localhost.csr