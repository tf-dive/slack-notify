FROM debian:bullseye-slim

RUN apt-get clean && apt-get update && apt-get upgrade -y \
    && apt-get install -y ca-certificates

# Non-root user `app`
RUN useradd app
WORKDIR /home/app

COPY bin/slack-notify ./

ENV LOGGING_LEVEL WARNING

COPY docker-entrypoint.sh ./

RUN chown -R app:app /home/app

# Change to user `app`
USER app

ENTRYPOINT ["./docker-entrypoint.sh"]

CMD ["./slack-notify"]

