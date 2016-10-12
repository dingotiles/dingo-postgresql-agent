FROM alpine:3.4

ENV PG_VERSION=9.6

# PostgreSQL
RUN set -x \
    && echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
    && apk update && apk add curl "postgresql@edge>9.6" "postgresql-contrib@edge>9.6" "postgresql-dev@edge>9.6"

ENV POSTGIS_VERSION 2.2.2
# https://hub.docker.com/r/unocha/alpine-base-postgis/~/dockerfile/
RUN echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    apk add --update --no-cache \
        perl \
        file \
        geos@testing \
        geos-dev@testing \
        libxml2-dev \
        gdal@testing \
        gdal-dev@testing \
        proj4@testing \
        proj4-dev@testing \
        gcc \
        make \
        libgcc g++ && \
    cd /tmp && \
    wget http://download.osgeo.org/postgis/source/postgis-${POSTGIS_VERSION}.tar.gz -O - | tar -xz && \
    chown root:root -R postgis-${POSTGIS_VERSION} && \
    cd /tmp/postgis-${POSTGIS_VERSION} && \
    ./configure && \
    echo "PERL = /usr/bin/perl" >> extensions/postgis/Makefile && \
    echo "PERL = /usr/bin/perl" >> extensions/postgis_topology/Makefile && \
    make -s && \
    make -s install  && \
    cd / && \
    rm -rf /tmp/postgis-${POSTGIS_VERSION} && \
    apk del perl file geos-dev \
            libxml2-dev gdal-dev proj4-dev \
            gcc make libgcc g++

# python3, wal-3
ENV WALE_VERSION=1.0.0
RUN set -x \
      && apk add --update bash python3 python3-dev jq curl alpine-sdk linux-headers musl-dev \
      && python3 -m ensurepip \
      && rm -r /usr/lib/python*/ensurepip \
      && pip3 install --upgrade pip setuptools \
      && rm -rf /root/.cache \
      && rm -rf /var/cache/apk/* \
      && pip install --no-cache-dir wal-e==${WALE_VERSION} awscli envdir pgxnclient --upgrade

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

ENV PATRONI_VERSION=1.1
RUN set -x \
    && curl -sL https://github.com/zalando/patroni/archive/v${PATRONI_VERSION}.tar.gz -o /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && tar xfz /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && mv /patroni-${PATRONI_VERSION} /patroni \
    && ls -al /patroni \
    && pip install --no-cache-dir -r /patroni/requirements.txt --upgrade \
    && rm -rf /root/.cache
#    && echo "pytz" >> /patroni/requirements.txt \

# 5432: PostgreSQL server
# 8008: Patroni API
EXPOSE 5432 8008

# Expose our data directory
VOLUME ["/data"]
ENV DATA_VOLUME=/data

CMD /bin/true
