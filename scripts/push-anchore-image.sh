#!/bin/bash

# This script will push images

set -e
set -u

docker push vilicus/anchore:latest
docker push vilicus/anchoredb:latest