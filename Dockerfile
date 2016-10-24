FROM dingotiles/dingo-postgresql96-agent-base:latest

ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN apk add --no-cache go git

# busybox-suid is to allow privilege escalation from a non-root user for setting up crontab for postgres
# RUN apk add --no-cache busybox-suid

# slug generates GUIDs
RUN set -x \
      && go get github.com/taskcluster/slugid-go/slug \
      && rm -rf $GOPATH/src

RUN set -x \
      && echo "@edge http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
      && apk add --no-cache --update pstree@edge

COPY config/patroni-default-values.yml /patroni/patroni-default-values.yml
COPY images/scripts/* /scripts/
COPY images/motd /etc/motd
RUN echo "source /etc/motd" >> /root/.bashrc
RUN echo "[[ -f /etc/patroni.d/.envrc ]] && source /etc/patroni.d/.envrc" >> /root/.bashrc
COPY images/supervisord.conf /etc/supervisor/supervisord.conf
COPY images/services/*.conf /etc/supervisor/conf.d/

COPY . /go/src/github.com/dingotiles/dingo-postgresql-agent
RUN set -x \
    && go install github.com/dingotiles/dingo-postgresql-agent \
    && rm -rf $GOPATH/src

ENV DINGO_API_URI https://api.dingotiles.com
CMD ["/scripts/entry.sh"]
env DINGO_IMAGE_VERSION=0.0.1
