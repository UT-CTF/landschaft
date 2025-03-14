#!/bin/bash

# Install deps
apt-get update
apt-get install -y pwgen docker.io docker-compose-v2

# Init /opt/graylog
mkdir -p /opt/graylog
mkdir -p /opt/graylog/tls_certs
cp docker-compose.yml /opt/graylog
cp $1 /opt/graylog/tls_certs/public-chain.pem
cp $2 /opt/graylog/tls_certs/private.key
# TODO: stop doing stupid things. (it's docker's fault)
chmod 666 /opt/graylog/tls_certs/*

# Firewall
ufw allow 12201/tcp
ufw allow 5555/tcp

# gen .env file
GRAYLOG_PASSWORD_SECRET=$(pwgen -N 1 -s 96)
echo "GRAYLOG_PASSWORD_SECRET=$GRAYLOG_PASSWORD_SECRET" > /opt/graylog/.env
GRAYLOG_ROOT_PASSWORD=$(pwgen -N 1 -s 96)
GRAYLOG_ROOT_PASSWORD_SHA2=$(echo -n $GRAYLOG_ROOT_PASSWORD | shasum -a 256)
echo "Graylog root password: $GRAYLOG_ROOT_PASSWORD"
echo "GRAYLOG_ROOT_PASSWORD_SHA2=$GRAYLOG_ROOT_PASSWORD_SHA2" >> /opt/graylog/.env

# sysctl
sysctl -w vm.max_map_count=262144
echo 'vm.max_map_count=262144' >> /etc/sysctl.conf

# Finish instructions
echo "1. To finish Graylog installation, run the following commands:"
echo "cd /opt/graylog"
echo "docker compose up -d"
echo "2. Graylog should be running on 127.0.0.1:9000!"
echo "3. Find the initial password in the logs via 'docker compose logs graylog'"
echo "4. Upload the ca-bundle.key file to the setup wizard"
echo "5. Login using admin:<root passswrod from above>"
echo "6. Upload the content pack"
echo "7. Setup TLS in the inputs settings"
