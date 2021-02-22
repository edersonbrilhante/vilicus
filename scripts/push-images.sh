#!/bin/bash

# This script will push images

set -e
set -u

docker push vilicus/clair:latest
docker push vilicus/clairdb:latest
docker push vilicus/trivy:latest
docker push vilicus/trivydb:latest
docker push vilicus/anchore:latest
docker push vilicus/anchoredb:latest
docker push vilicus/vilicus:latest