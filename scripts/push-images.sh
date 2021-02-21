#!/bin/bash

# This script will push all images

set -e
set -u

docker push "vilicus/anchore:latest"
docker push "vilicus/clair:latest"
docker push "vilicus/vilicus:latest"
docker push "vilicus/trivy:base"
docker push "vilicus/postgres:base"
