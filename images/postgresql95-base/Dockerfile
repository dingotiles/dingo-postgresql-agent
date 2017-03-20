FROM dingotiles/dingo-postgresql-base:latest


# PostgreSQL
RUN set -x \
    && apk update && apk add postgresql postgresql-contrib postgresql-dev

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

ENV PATRONI_VERSION=1.2.3
RUN set -x \
    && curl -sL https://github.com/zalando/patroni/archive/v${PATRONI_VERSION}.tar.gz -o /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && tar xfz /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && mv /patroni-${PATRONI_VERSION} /patroni \
    && pip3 install --no-cache-dir -r /patroni/requirements.txt --upgrade \
    && rm -rf /root/.cache /tmp/*
