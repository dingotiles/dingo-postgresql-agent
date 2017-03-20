FROM dingotiles/dingo-postgresql-base:latest

RUN apk add --update --virtual=.build-dependencies perl bison flex readline zlib python3-dev

ENV POSTGRESQL_VERSION=9.5.6 PG_VERSION=9.5.6
RUN wget https://ftp.postgresql.org/pub/source/v${PG_VERSION}/postgresql-${PG_VERSION}.tar.gz -O /tmp/postgresql-${PG_VERSION}.tar.gz && \
      echo "40558693d9a2c30d5e389a5436812458  /tmp/postgresql-${PG_VERSION}.tar.gz" | md5sum -c - && \
      tar xfz /tmp/postgresql-${PG_VERSION}.tar.gz -C /tmp/ && \
      cd /tmp/postgresql-${PG_VERSION} && \
      ./configure --prefix=/usr/local \
        --with-python \
        --without-openssl \
        --without-readline --without-zlib \
        PYTHON=/usr/bin/python3 && \
      make && make install && \
      rm -rf /tmp/postgresql-${PG_VERSION}*

ENV PATRONI_VERSION=1.2.3
RUN set -x \
    && curl -sL https://github.com/zalando/patroni/archive/v${PATRONI_VERSION}.tar.gz -o /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && tar xfz /tmp/patroni-${PATRONI_VERSION}.tar.gz \
    && mv /patroni-${PATRONI_VERSION} /patroni \
    && pip3 install --no-cache-dir -r /patroni/requirements.txt --upgrade \
    && rm -rf /root/.cache /tmp/*
