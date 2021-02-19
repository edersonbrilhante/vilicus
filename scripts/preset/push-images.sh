#!/bin/bash

# This script will push postgres preset images

set -e
set -u

docker push vilicus/postgres:preset-files
docker push vilicus/postgres:preset-volume