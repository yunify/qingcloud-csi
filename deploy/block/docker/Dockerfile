FROM ubuntu:16.04
MAINTAINER "Yunify"

RUN apt-get update && \
    apt-get install -y ca-certificates ca-certificates-java
COPY blockplugin /blockplugin
RUN chmod +x /blockplugin && \
    mkdir -p /var/log/blockplugin

ENTRYPOINT ["/blockplugin"]
