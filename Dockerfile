FROM alpine
LABEL maintainer="Adel Abdelhak"

ENV LISTEN_ADDR=:9191 TELEMETRY_PATH=/metrics CB_URI=http://localhost:8091 CB_ADMIN_USER=admin CB_ADMIN_PASSWORD=password

ADD ./dist/couchbase_exporter /bin/couchbase_exporter

CMD /bin/couchbase_exporter