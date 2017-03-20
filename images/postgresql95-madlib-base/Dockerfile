FROM dingotiles/dingo-postgresql-base:latest

RUN apk add --update --virtual=.build-dependencies perl bison flex readline zlib python python-dev

ENV POSTGRESQL_VERSION=9.5.6 PG_VERSION=9.5.6
RUN wget https://ftp.postgresql.org/pub/source/v${PG_VERSION}/postgresql-${PG_VERSION}.tar.gz -O /tmp/postgresql-${PG_VERSION}.tar.gz && \
      echo "40558693d9a2c30d5e389a5436812458  /tmp/postgresql-${PG_VERSION}.tar.gz" | md5sum -c - && \
      tar xfz /tmp/postgresql-${PG_VERSION}.tar.gz -C /tmp/ && \
      cd /tmp/postgresql-${PG_VERSION} && \
      ./configure --prefix=/usr/local \
        --with-python \
        --without-openssl \
        --without-readline --without-zlib \
        PYTHON=/usr/bin/python && \
      make && make install && \
      rm -rf /tmp/postgresql-${PG_VERSION}*

ENV PATRONI_VERSION=1.2.3
RUN set -x \
    && curl -sL https://github.com/zalando/patroni/archive/v${PATRONI_VERSION}.tar.gz -o /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && tar xfz /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && mv /patroni-${PATRONI_VERSION} /patroni \
    && pip3 install --no-cache-dir -r /patroni/requirements.txt --upgrade \
    && rm -rf /root/.cache /tmp/*

RUN pip install pgxnclient
RUN echo "@testing http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories && \
    echo "http://dl-2.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "http://dl-3.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "http://dl-5.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    apk add --update --no-cache \
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
      pv geos proj4 gdal \
      m4 \
      build-base gcc abuild binutils binutils-doc gcc-doc \
      cmake cmake-doc
RUN pgxn install madlib
