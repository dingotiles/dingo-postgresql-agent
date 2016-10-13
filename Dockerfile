FROM dingotiles/dingo-postgresql96-agent-base:latest

ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN apk add --no-cache go git

COPY . /go/src/github.com/dingotiles/dingo-postgresql-agent
RUN set -x\
    && go install github.com/dingotiles/dingo-postgresql-agent

COPY config/patroni-default-values.yml /patroni/patroni-default-values.yml
COPY images/scripts/* /scripts/
COPY images/supervisord.conf /etc/supervisor/supervisord.conf
COPY images/services/*.conf /etc/supervisor/conf.d/
