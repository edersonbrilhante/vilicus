#!/bin/bash

# This script will build vilicus image

set -e
set -u

docker build -f deployments/dockerfiles/vilicus/Dockerfile -t vilicus/vilicus:latest .