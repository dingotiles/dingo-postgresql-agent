FROM alpine:3.4

ENV PG_VERSION=9.6

# python3, psql, wal-3
RUN echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
      apk update && apk add curl "postgresql@edge>9.6" "postgresql-contrib@edge>9.6"

RUN apk add --update bash python3 python3-dev jq curl alpine-sdk linux-headers && \
      python3 -m ensurepip && \
      rm -r /usr/lib/python*/ensurepip && \
      pip3 install --upgrade pip setuptools && \
      rm -rf /root/.cache && \
      rm -rf /var/cache/apk/* && \
      pip install wal-e awscli envdir --upgrade

# etcdctl
RUN curl -sL https://github.com/coreos/etcd/releases/download/v2.3.4/etcd-v2.3.4-linux-amd64.tar.gz -o /tmp/etcd-v2.3.4-linux-amd64.tar.gz && \
      tar xzvf /tmp/etcd-v2.3.4-linux-amd64.tar.gz -C /tmp && \
      mv /tmp/etcd-v2.3.4-linux-amd64/etcdctl /usr/local/bin && \
      rm -rf /tmp/etcd*

CMD /bin/true
