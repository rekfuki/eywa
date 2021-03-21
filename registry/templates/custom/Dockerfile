FROM registry.eywa.rekfuki.dev/of-watchdog:latest as watchdog
FROM alpine:3.12
# Add non root user and certs
RUN apk --no-cache add ca-certificates tree \
    && addgroup -S app && adduser -S -g app app \
    && mkdir -p /home/app \
    && chown app /home/app

WORKDIR /home/app

ARG FPROCESS
COPY --chown=app ./function/  ./
COPY --from=watchdog --chown=app /of-watchdog ./fwatchdog
RUN chmod +x ./fwatchdog

USER app

RUN pwd
RUN ls -la
RUN env 

RUN if [[ ! -e $FPROCESS ]]; then >&2 echo "Provided executable ${FPROCESS} does not exist..."; exit 1; fi
RUN if output=$(ldd "$FPROCESS" 2>&1); then true; elif echo "$output" | grep 'Not a valid dynamic program' >/dev/null; then true; else echo "$output"; echo "ERROR: Custom runtimes must be compilled statically"; false; fi

RUN chmod +x $FPROCESS

RUN ls -la
ENV fprocess="./${FPROCESS}"
ENV mode="http"
ENV upstream_url="http://127.0.0.1:8082"
RUN env

CMD ["./fwatchdog"]