FROM dingotiles/dingo-postgresql:pre-test

COPY wal-e /tmp/wal-e

RUN cd /tmp/wal-e && pip3 install .

ENV WALE_VERSION=$(wal-e version)
