FROM alpine:3.4

ENV PG_VERSION=9.6

# PostgreSQL
RUN set -x \
    && echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
    && apk update && apk add curl "postgresql@edge>9.6" "postgresql-contrib@edge>9.6"

# python3, wal-3
RUN set -x \
      && apk add --update bash python3 python3-dev jq curl alpine-sdk linux-headers \
      && python3 -m ensurepip \
      && rm -r /usr/lib/python*/ensurepip \
      && pip3 install --upgrade pip setuptools \
      && rm -rf /root/.cache \
      && rm -rf /var/cache/apk/* \
      && pip install wal-e awscli envdir --upgrade

# etcdctl
RUN set -x \
      && curl -sL https://github.com/coreos/etcd/releases/download/v2.3.4/etcd-v2.3.4-linux-amd64.tar.gz -o /tmp/etcd-v2.3.4-linux-amd64.tar.gz \
      && tar xzvf /tmp/etcd-v2.3.4-linux-amd64.tar.gz -C /tmp \
      && mv /tmp/etcd-v2.3.4-linux-amd64/etcdctl /usr/local/bin \
      && rm -rf /tmp/etcd*

# gosu https://github.com/tianon/gosu
# https://github.com/tianon/gosu/releases
ENV GOSU_VERSION=1.10
RUN set -x \
    && apk add --no-cache --update dpkg gnupg openssl \
    && dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')" \
    && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch" \
    && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc" \
    && export GNUPGHOME="$(mktemp -d)" \
    && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
    && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
    && rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
    && chmod +x /usr/local/bin/gosu \
    && gosu nobody true \
    && rm -rf /var/cache/apk/*

# 5432: PostgreSQL server
# 8008: Patroni API
EXPOSE 5432 8008

# Expose our data directory
VOLUME ["/data"]
ENV DATA_VOLUME=/data

CMD /bin/true
