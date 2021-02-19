#!/bin/bash

# This script will push postgres preset images

set -e
set -u



docker tag "vilicus/postgres:preset-files-latest" "vilicus/postgres:preset-files"
docker tag "vilicus/postgres:preset-volume-latest" "vilicus/postgres:preset-volume"

docker push vilicus/postgres:preset-files
docker push vilicus/postgres:preset-volume