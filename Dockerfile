FROM dingotiles/dingo-postgresql95-agent-base:latest

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

COPY . /go/src/github.com/dingotiles/dingo-postgresql-agent
RUN set -x \
    && go install github.com/dingotiles/dingo-postgresql-agent \
    && rm -rf $GOPATH/src
