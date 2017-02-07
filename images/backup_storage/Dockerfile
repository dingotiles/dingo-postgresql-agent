FROM ubuntu:16.04

RUN apt-get update && apt-get install -y openssh-server ssh vim
RUN mkdir /var/run/sshd

EXPOSE 22
VOLUME /data
COPY scripts/* /scripts/

CMD ["/scripts/run.sh"]
