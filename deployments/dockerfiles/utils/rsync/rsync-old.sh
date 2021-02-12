#!/bin/bash

# Based in https://github.com/neam/docker-diff-based-layers

echo "[vilicus_conf]" > /tmp/.rsyncd.conf;
echo "path = ${VOLUME}" >> /tmp/.rsyncd.conf;
echo "uid = postgres" >> /tmp/.rsyncd.conf;
echo "gid = postgres" >> /tmp/.rsyncd.conf;
echo "read only = false" >> /tmp/.rsyncd.conf;

rsync --daemon --port 873 --no-detach -vv --config /tmp/.rsyncd.conf
