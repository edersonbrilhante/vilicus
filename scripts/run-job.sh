#!/bin/bash

# This script will run vilicus scan

set -e
set -u

function set_vars (){
  if [ -z ${TEMPLATE+x} ]; then
    echo '$TEMPLATE is unset'
    exit 2
  fi
  
  if [ -z ${OUTPUT+x} ]; then
    echo '$TEMPLATE is unset'
    exit 2
  fi
  
  if [ -z ${IMAGE+x} ]; then
    echo '$IMAGE is unset'
    exit 2
  fi
  
  if [ -z ${CONFIG+x} ]; then
    CONFIG=/opt/vilicus/configs/conf.yaml
  fi

  CMD='dockerize -wait http://vilicus:8080/healthz -wait-retry-interval 60s -timeout 2000s vilicus-client'
}

function push_image (){
  echo "Push Image"
  if [[ $IMAGE =~ "localhost:5000" ]]; then
    docker push $IMAGE
    IMAGE="${IMAGE/localhost:5000/localregistry.vilicus.svc:5000}" 
  fi
}

function download_docker_compose (){
  echo "Download Docker Compose"
  wget -O docker-compose.yml https://raw.githubusercontent.com/edersonbrilhante/vilicus/main/deployments/docker-compose.yml 
}

function run_docker_compose (){
  echo "Run Docker Compose"
  docker-compose up -d
}

function create_artifacts (){
  echo "Create Artifacts Folder"
  mkdir -p artifacts && chmod 777 artifacts
}

function run_scan (){
  echo "Run Scan"
  docker run \
    -v ${PWD}/artifacts:/artifacts \
    --network container:vilicus \
    vilicus/vilicus:latest \
    sh -c "${CMD} -p ${CONFIG} -i ${IMAGE}  -t ${TEMPLATE} -o ${OUTPUT}"
}

set_vars
download_docker_compose
run_docker_compose
create_artifacts
push_image
run_scan