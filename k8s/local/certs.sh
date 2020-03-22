#!/bin/bash

echo "Creating self-signed CA certificates for TLS and installing them in the local trust stores"
CA_CERTS_FOLDER=$(pwd)/.certs # This requires mkcert to be installed/available
ENVIRONMENT_DEV=dev
echo ${CA_CERTS_FOLDER}
rm -rf ${CA_CERTS_FOLDER}
mkdir -p ${CA_CERTS_FOLDER}
mkdir -p ${CA_CERTS_FOLDER}/${ENVIRONMENT_DEV} # The CAROOT env variable is used by mkcert to determine where to read/write files# Reference: https://github.com/FiloSottile/mkcert
export CAROOT=${CA_CERTS_FOLDER}/${ENVIRONMENT_DEV}
mkcert -install

echo "Creating K8S secrets with the CA private keys (will be used by the cert-manager CA Issuer)"
kubectl -n cert-manager create secret tls local-https-secret --key=${CA_CERTS_FOLDER}/${ENVIRONMENT_DEV}/rootCA-key.pem --cert=${CA_CERTS_FOLDER}/${ENVIRONMENT_DEV}/rootCA.pem