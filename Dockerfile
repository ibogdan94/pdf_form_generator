FROM golang:1.11

RUN mkdir -p $GOPATH/bin && \
    go get github.com/cortesi/modd/cmd/modd

RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

RUN apt-get -q -y install libjpeg-dev libpng-dev libtiff-dev \
    libgif-dev libx11-dev --no-install-recommends \
    ghostscript

ENV IMAGEMAGICK_VERSION=7.0.8-12

RUN cd && \
	wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
	tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
	cd ImageMagick* && \
	./configure \
	    --without-magick-plus-plus \
	    --without-perl \
	    --disable-openmp \
	    --with-gvc=no \
	    --disable-docs && \
	make -j$(nproc) && make install && \
	ldconfig /usr/local/lib


ADD . /go/src/pdf_form_generator
WORKDIR /go/src/pdf_form_generator

ENV GO111MODULE=on
RUN go get

EXPOSE 8080

CMD modd