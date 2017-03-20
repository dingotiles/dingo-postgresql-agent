FROM dingotiles/dingo-postgresql:pre-test

RUN pip install pgxnclient
RUN echo "http://dl-2.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
    echo "http://dl-3.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
    echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
    echo "http://dl-5.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories
RUN apk add --update --no-cache pv geos proj4 gdal \
      m4 \
      build-base gcc abuild binutils binutils-doc gcc-doc \
      cmake cmake-doc
RUN pgxn install madlib
ENV DINGO_IMAGE_NAME=dingotiles/dingo-postgresql-madlib
