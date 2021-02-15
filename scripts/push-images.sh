#!/bin/bash

# This script will push all images

anchoreVersion=$(docker image inspect vilicus/anchore:latest --format '{{ index .Config.Labels "vilicus.app.version"}}')
clairVersion=$(docker image inspect vilicus/clair:latest --format '{{ index .Config.Labels "vilicus.app.version"}}')
trivyVersion=$(docker image inspect vilicus/trivy:latest --format '{{ index .Config.Labels "vilicus.app.version"}}')
vilicusVersion=$(docker image inspect vilicus/vilicus:latest --format '{{ index .Config.Labels "vilicus.app.version"}}')

docker tag "vilicus/anchore:latest" "vilicus/anchore:$anchoreVersion"
docker tag "vilicus/clair:latest" "vilicus/clair:$clairVersion"
docker tag "vilicus/trivy:latest" "vilicus/trivy:$trivyVersion"
docker tag "vilicus/vilicus:latest" "vilicus/vilicus:$vilicusVersion"

docker push "vilicus/anchore:latest" && docker push "vilicus/anchore:$anchoreVersion"
docker push "vilicus/clair:latest" && docker push "vilicus/clair:$clairVersion"
docker push "vilicus/trivy:latest" && docker push "vilicus/trivy:$trivyVersion"
docker push "vilicus/vilicus:latest" && docker push "vilicus/vilicus:$vilicusVersion"
docker push "vilicus/postgres:base"
