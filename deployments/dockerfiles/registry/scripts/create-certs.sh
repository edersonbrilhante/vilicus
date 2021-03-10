#!/bin/bash

openssl req -new \
    -newkey rsa:2048 \
    -nodes \
    -subj "/C=GB/ST=England/L=Brighton/O=Example/CN=*.vilicus.svc/emailAddress=info@vilicus.svc" \
    -reqexts SAN \
    -config <(cat /etc/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:localhost,DNS:localregistry.vilicus.svc")) \
    -keyout /opt/vilicus/certs/vilicus.key \
    -x509 \
    -days 3650 \
    -extensions SAN \
    -out /opt/vilicus/certs/vilicus.crt