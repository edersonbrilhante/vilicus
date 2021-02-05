#!/bin/bash

# This script will build all images

docker build -f deployments/dockerfiles/anchore/Dockerfile -t vilicus/anchore:latest .
docker build -f deployments/dockerfiles/clair/Dockerfile -t vilicus/clair:latest .
docker build -f deployments/dockerfiles/trivy/Dockerfile -t vilicus/trivy:latest .
docker build -f deployments/dockerfiles/vilicus/Dockerfile -t vilicus/vilicus:latest .