#!/bin/bash

# This script will push images

set -e
set -u

docker push vilicus/trivy:latest
docker push vilicus/trivydb:latest