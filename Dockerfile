FROM debian:jessie-slim
RUN apt-get clean && apt-get update && apt-get upgrade -y \
    && apt-get install -y ca-certificates

COPY rootfs /

CMD /slack-notify
