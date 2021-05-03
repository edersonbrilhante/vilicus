# Vilicus

<p align="left">
  <a href="https://github.com/edersonbrilhante/vilicus/releases"><img src="https://img.shields.io/github/v/release/edersonbrilhante/vilicus"/></a>
  <a href="https://travis-ci.com/edersonbrilhante/vilicus.svg?branch=main"><img src="https://travis-ci.com/edersonbrilhante/vilicus.svg?branch=main"/></a>
</p>

# Table of Contents
- [Overview](#overview)
  - [How does it work?](#how-does-it-work)
- [Architecture](#architecture)
- [Development](#development)
    - [Run deployment manually](#run-deployment-manually)
- [Usage](#usage)
    - [Example of analysis](#example-of-analysis)

---

## Overview
Vilicus is an open source tool that orchestrates security scans of container images(docker/oci) and centralizes all results into a database for further analysis and metrics. It can perform using [Anchore](https://github.com/anchore/anchore-engine), [Clair](https://github.com/quay/clair) and [Trivy](https://github.com/aquasecurity/trivy).

### How does it work?
There many tools to scan container images, but sometimes the results can be diferent in each one them. So the main goal of this project is to help development teams improve the quality of their container images by finding vulnerabilities and thus addressing them with anagnostic sight from vendors.

**Here you can find articles comparing the scanning tools**:
- [Open Source CVE Scanner Round-Up: Clair vs Anchore vs Trivy](https://boxboat.com/2020/04/24/image-scanning-tech-compared/)
- [5 open source tools for container security](https://opensource.com/article/18/8/tools-container-security)

---

## Architecture
![Kiku](docs/arch.gif)

---

## Development
### Run deployment manually
```bash
docker-compose -f deployments/docker-compose.yaml up -d
```

---

## Usage

### Requirements
- Disk Space ~30GB:
  - Docker System:
    - Images ~14GB
    - Containers ~11GB
    - Local Volumes ~200MB
- Docker
- Docker Compose
- Bash
- Wget

### Using vilicus client
Run these following commands:
```
export TEMPLATE=<template>
export OUTPUT=<output>
export IMAGE=<public_image>|<vilicus_local_image>
wget -O run-job.sh https://raw.githubusercontent.com/edersonbrilhante/vilicus/main/scripts/run-job.sh
chmod +x ./run-job.sh
./run-job.sh
```
The result will be stored in into the file set by the environment variable `OUTPUT`.

### Templates and Outputs
**Gitlab**<br>
***Template***: `/opt/vilicus/contrib/gitlab.tpl`<br>
***Output***: `/artifacts/gl-container-scanning-report.json`

**Sarif**<br>
***Template***: `/opt/vilicus/contrib/sarif.tpl`<br>
***Output***: `/artifacts/result.sarif`

### Public image and Local images
Vilicus provides support images hosted in public repository and local builds. Public image is an image hosted in public repository such as DockerHub. To scan images in self-hosted registry or local build you must tag the image to the vilicus local registry.

**Self-hosted registry**
`docker tag <self-hosted-registry>/<image:tag> localhost:5000/<image:tag>`

**Local build**
`docker build -t localhost:5000/<image:tag> -f <Dockerfile> <context>`

### Free Online Service
Vilicus also provides a [free online service](http://vilicus.edersonbrilhante.com.br/). 

#### How it works?
This service is a serverless full-stack application with backend workers and database only using git and ci/cd runners.

The Frontend is hosted in GitHub Pages. This frontend is a landing page with a free service to scan or display the vulnerabilities in container images.

The results of container image scans are stored in a GitLab Repository.

When the user asks to show the results from an image, the frontend consumes the GitLab API to retrieve the file with vulns from this image. In case this image is not scanned yet, the user has the option to schedule a scan using a google form.

When this form is filled, the data is sent to a Google Spreadsheet.

A GitHub Workflow runs every 5 minutes to check if there are new answers in this Spreadsheet. For each new image in the Spreadsheet, this workflow triggers another Workflow to scan the image and save the result in the GitLab Repository.