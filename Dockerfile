FROM dingotiles/dingo-postgresql96-agent-base:latest

# Set up GOPATH
ADD . /go/src/github.com/dingotiles/dingo-postgresql-agent
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN set -x \
    && apk add --update --no-cache go git

RUN set -x\
    && go install github.com/dingotiles/dingo-postgresql-agent

ADD config/patroni-default-values.yml /patroni/patroni-default-values.yml
ADD images/scripts/* /scripts/
ADD images/supervisord.conf /etc/supervisor/supervisord.conf
ADD images/services/*.conf /etc/supervisor/conf.d/
ADD sample-patroni.yml /config/patroni.yml
