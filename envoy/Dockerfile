FROM alpine:latest AS build

ENV NGINX_VERSION 1.19.8
ENV DEVEL_KIT_VERSION 0.3.1
ENV SET_MISC_VERSION 0.32
ENV HEADERS_MORE_VERSION 0.33

# Prerequisites
RUN apk add build-base pcre-dev openssl-dev zlib-dev

# Fetch and unpack sources
ADD https://nginx.org/download/nginx-$NGINX_VERSION.tar.gz /build/nginx.tar.gz
ADD https://github.com/vision5/ngx_devel_kit/archive/v$DEVEL_KIT_VERSION.tar.gz /build/devel-kit.tar.gz
ADD https://github.com/openresty/set-misc-nginx-module/archive/v$SET_MISC_VERSION.tar.gz /build/set-misc.tar.gz
ADD https://github.com/openresty/headers-more-nginx-module/archive/v$HEADERS_MORE_VERSION.tar.gz /build/headers-more.tar.gz

RUN mkdir /build/nginx     && tar xf /build/nginx.tar.gz     -C /build/nginx     --strip-components=1 && ls /build/nginx
RUN mkdir /build/devel-kit && tar xf /build/devel-kit.tar.gz -C /build/devel-kit --strip-components=1 && ls /build/devel-kit
RUN mkdir /build/set-misc  && tar xf /build/set-misc.tar.gz  -C /build/set-misc  --strip-components=1 && ls /build/set-misc
RUN mkdir /build/headers-more  && tar xf /build/headers-more.tar.gz  -C /build/headers-more  --strip-components=1 && ls /build/headers-more

WORKDIR /build/nginx

# Compile and install
RUN ./configure \
    --prefix=/usr \
    --with-http_auth_request_module \
    --with-http_ssl_module \
    --with-http_stub_status_module \
    --add-module=/build/devel-kit \
    --add-module=/build/set-misc \
    --add-module=/build/headers-more \
    && make -j4 \
    && make install

###

FROM alpine:latest

# Prerequisites
RUN apk --no-cache add pcre openssl zlib

# Install
COPY --from=build /usr/sbin/nginx /usr/sbin
RUN mkdir -p /usr/logs

ADD start.sh reload.sh watcher /
ADD conf /conf

ENTRYPOINT ["/start.sh"]
