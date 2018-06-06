#!/usr/bin/env bash
# starts containers for cloudsqlproxy, and for geoserver
# then puts them on the same network so they can talk to one another

docker run --name="cloudsql" -d -v /Users/hdh/bin/cloud_sql_proxy.linux.amd64:/cloudsql \
-v ~/UNSAFE_PERSONAL/habitat-dev-client.json:/config -p 5432:5432 \
gcr.io/cloudsql-docker/gce-proxy:1.11 /cloudsql -instances=hdh-habitat-modeling:us-west1:gis=tcp:0.0.0.0:5432 -credential_file=/config

docker run --name "geoserver" -v ~/ziggy-data/geoserver-data:/opt/geoserver/data_dir -p 8080:8080 -d -t kartoza/geoserver

docker network connect geo-net cloudsql
docker network connect geo-net geoserver
