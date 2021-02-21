#!/bin/bash

# This script will build all images

set -e
set -u

docker build -f deployments/dockerfiles/anchore/Dockerfile -t vilicus/anchore:latest .
docker build -f deployments/dockerfiles/clair/Dockerfile -t vilicus/clair:latest .
docker build -f deployments/dockerfiles/trivy/Dockerfile -t vilicus/trivy:base .
docker build -f deployments/dockerfiles/vilicus/Dockerfile -t vilicus/vilicus:latest .
docker build -f deployments/dockerfiles/postgres/Dockerfile -t vilicus/postgres:base .