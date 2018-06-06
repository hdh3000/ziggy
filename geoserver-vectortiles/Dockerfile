# This is a hacky way of using the kartoza/geoserver container, and then adding the vector tiles plug-in
# I say hacky because the version will probably break, and it's a local copy of the plugin

FROM kartoza/geoserver

COPY ./geoserver-2.13.1-vectortiles-plugin /usr/local/tomcat/webapps/geoserver/WEB-INF/lib/