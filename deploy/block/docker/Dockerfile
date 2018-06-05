FROM ubuntu:16.04
MAINTAINER "Yunify"

COPY blockplugin /blockplugin

RUN apt-get update && \
    apt-get install -y ca-certificates ca-certificates-java
RUN chmod +x /blockplugin && \
    mkdir -p /var/log/blockplugin

ENTRYPOINT ["/blockplugin"]
