FROM dingotiles/dingo-postgresql95-agent-base:latest

# slug generates GUIDs
RUN set -x \
      && go get github.com/taskcluster/slugid-go/slug \
      && rm -rf $GOPATH/src

RUN set -x \
      && echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
      && apk add --no-cache --update pstree@edge

COPY config/patroni-default-values.yml /patroni/patroni-default-values.yml
COPY images/motd /etc/motd
RUN echo "source /etc/motd" >> /root/.bashrc
RUN echo "[[ -f /etc/patroni.d/.envrc ]] && source /etc/patroni.d/.envrc" >> /root/.bashrc
COPY images/supervisord.conf /etc/supervisor/supervisord.conf
COPY images/services/*.conf /etc/supervisor/conf.d/

COPY images/scripts /scripts
CMD ["/scripts/entry.sh"]

ENV DINGO_API_URI=https://api.dingotiles.com DINGO_IMAGE_VERSION=0.0.0
VOLUME ["/backups"]
ENV LOCAL_BACKUP_VOLUME=/backups/

COPY . /go/src/github.com/dingotiles/dingo-postgresql-agent
RUN set -x \
    && go install github.com/dingotiles/dingo-postgresql-agent \
    && rm -rf $GOPATH/src

COPY wal-e /tmp/wal-e

RUN cd /tmp/wal-e && pip3 install . --upgrade
