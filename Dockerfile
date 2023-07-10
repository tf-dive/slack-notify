FROM debian:bullseye-slim
RUN apt-get clean && apt-get update && apt-get upgrade -y \
    && apt-get install -y ca-certificates

COPY rootfs /

CMD /slack-notify
