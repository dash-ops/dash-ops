FROM alpine

RUN apk add --update ca-certificates && \
  rm -rf /var/cache/apk/* /tmp/*
EXPOSE 8080

COPY dash-ops /
COPY front/dist /app

CMD ["/dash-ops"]
