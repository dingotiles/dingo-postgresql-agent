FROM dingotiles/dingo-postgresql95-src-agent-base:latest

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
ENV DINGO_IMAGE_NAME=dingotiles/dingo-postgresql-madlib
