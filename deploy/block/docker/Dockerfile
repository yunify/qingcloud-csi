FROM ubuntu:16.04
MAINTAINER "Yunify"

COPY blockplugin /blockplugin

RUN chmod +x /blockplugin && \
    mkdir -p /var/log/blockplugin

ENTRYPOINT ["/blockplugin"]
