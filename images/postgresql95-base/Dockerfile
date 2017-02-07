FROM alpine:3.4
# contains postgresql-9.5.4

ENV PG_VERSION=9.5

ENV GOLANG_VERSION 1.7.4
ENV GOLANG_SRC_SHA256 4c189111e9ba651a2bb3ee868aa881fab36b2f2da3409e80885ca758a6b614cc
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		bash gcc musl-dev openssl go \
	\
	&& export GOROOT_BOOTSTRAP="$(go env GOROOT)" \
	\
	&& wget -q "https://golang.org/dl/go${GOLANG_VERSION}.src.tar.gz" -O golang.tar.gz \
	&& echo "$GOLANG_SRC_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz \
	&& cd /usr/local/go/src \
	&& ./make.bash \
	\
	&& rm -rf /*.patch \
	&& apk del .build-deps \
  && mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# PostgreSQL
RUN set -x \
    && apk update && apk add bash curl openssh-client postgresql postgresql-contrib postgresql-dev

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
ENV WALE_VERSION=1.0.2
RUN set -x \
      && echo "@community http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
      && apk add --update python3 python3-dev bash sed jq curl alpine-sdk linux-headers musl-dev pv@community \
      && python3 -m ensurepip \
      && rm -r /usr/lib/python*/ensurepip \
      && pip3 install --upgrade pip setuptools \
      && rm -rf /root/.cache \
      && rm -rf /var/cache/apk/* \
      && pip3 install --no-cache-dir wal-e[aws,azure,google,swift]==${WALE_VERSION} awscli envdir pgxnclient --upgrade

# python2, dumb-init, supervisor
ENV SUPERVISOR_VERSION=3.3.1
RUN set -x \
      && apk add --update bash python py-pip \
      && rm -rf /var/cache/apk/* \
      && pip install --upgrade pip setuptools \
      && pip install --no-cache-dir supervisor==${SUPERVISOR_VERSION} dumb-init --upgrade

# etcdctl
ENV ETCD_VERSION=2.3.7
RUN set -x \
      && curl -sL https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz -o /tmp/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz \
      && tar xzvf /tmp/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz -C /tmp \
      && mv /tmp/etcd-v${ETCD_VERSION}-linux-amd64/etcdctl /usr/local/bin \
      && rm -rf /tmp/etcd*

ENV PATRONI_VERSION=1.2.3
RUN set -x \
    && curl -sL https://github.com/zalando/patroni/archive/v${PATRONI_VERSION}.tar.gz -o /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && tar xfz /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && mv /patroni-${PATRONI_VERSION} /patroni \
    && pip3 install --no-cache-dir -r /patroni/requirements.txt --upgrade \
    && rm -rf /root/.cache /tmp/*

ENV SPRUCE_VERSION=1.8.7
RUN set -x \
    && curl -L https://github.com/geofffranks/spruce/releases/download/v${SPRUCE_VERSION}/spruce-linux-amd64 -o /usr/local/bin/spruce \
    && chmod +x /usr/local/bin/spruce

ENV S3CLI_VERSION=0.0.51
RUN set -x \
      && curl -L https://s3.amazonaws.com/s3cli-artifacts/s3cli-${S3CLI_VERSION}-linux-amd64 -o /usr/local/bin/s3cli \
      && chmod +x /usr/local/bin/s3cli

# 5432: PostgreSQL server
# 8008: Patroni API
EXPOSE 5432 8008

# Expose our data directory
VOLUME ["/data"]
ENV DATA_VOLUME=/data

CMD ["/bin/true"]
