#!/bin/bash

# This script will build registry image

set -e
set -u

docker build -f deployments/dockerfiles/registry/Dockerfile -t vilicus/registry:latest .