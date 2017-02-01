FROM ubuntu:16.04

RUN apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends rsync && \
  apt-get clean autoclean && \
  apt-get autoremove -y && \
  rm -rf /var/lib/apt/lists/*

EXPOSE 22
EXPOSE 873
VOLUME /data
COPY scripts/* /scripts/

ENTRYPOINT ["/scripts/run.sh"]
