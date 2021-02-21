#!/bin/bash

# This script will push preset images

set -e
set -u



docker tag "vilicus/postgres:preset-files-latest" "vilicus/postgres:preset-files"
docker tag "vilicus/postgres:preset-volume-latest" "vilicus/postgres:preset-volume"
docker tag "vilicus/trivy:preset-latest" "vilicus/trivy:preset"

docker push vilicus/postgres:preset-files
docker push vilicus/postgres:preset-volume
docker push vilicus/trivy:preset