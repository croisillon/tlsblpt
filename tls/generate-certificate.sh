#!/bin/bash

caCertDays=365
serverCertDays=60
clientCertDays=60

certDir="../certs"
caDir="${certDir}"/ca
serverDir="${certDir}"/server
clientDir="${certDir}"/client

C="C=RU"
ST="ST=Russia"
L="L=Chelyabinsk"
O="O=HomeOffice"
OU="OU=IT"
CN="CN=localhost"
emailAddress="emailAddress=root@localhost"

subj="/$C/$ST/$L/$O/$OU/$CN/$emailAddress"

mkdir -p $certDir
mkdir -p $caDir
mkdir -p $serverDir
mkdir -p $clientDir

# 1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days $caCertDays -nodes -keyout "${caDir}"/key.pem -out "${caDir}"/cert.pem -subj "${subj}"

# echo "CA's self-signed certificate"
openssl x509 -in "${caDir}"/key.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout "${serverDir}"/key.pem -out "${serverDir}"/req.pem -subj "${subj}"

# # 3. Use CA's private key to sign web server's CSR and get back the signed certificate
openssl x509 -req -in "${serverDir}"/req.pem -days ${serverCertDays} -CA "${caDir}"/cert.pem -CAkey "${caDir}"/key.pem -CAcreateserial -out "${serverDir}"/cert.pem -extfile "self-signed-cert.ext"

# echo "Server's signed certificate"
openssl x509 -in "${serverDir}"/cert.pem -noout -text

# # 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout "${clientDir}"/key.pem -out "${clientDir}"/req.pem -subj "${subj}"

# # 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in "${clientDir}"/req.pem -days ${clientCertDays} -CA "${caDir}"/cert.pem -CAkey "${caDir}"/key.pem -CAcreateserial -out "${clientDir}"/cert.pem -extfile "self-signed-cert.ext"

# echo "Client's signed certificate"
openssl x509 -in "${clientDir}"/cert.pem -noout -text