#!/usr/bin/env bash
docker stop cloudsql
docker rm cloudsql

docker stop geoserver
docker rm geoserver