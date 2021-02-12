#!/bin/bash

# This script will push postgres images

docker tag vilicus/postgres-presets:processed vilicus/postgres-presets:latest
docker tag vilicus/postgres-presets:processed vilicus/postgres-presets:monthly

docker push vilicus/postgres-presets:monthly
docker push vilicus/postgres-presets:latest