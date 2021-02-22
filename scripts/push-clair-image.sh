#!/bin/bash

# This script will push images

set -e
set -u

docker push vilicus/clair:latest
docker push vilicus/clairdb:latest