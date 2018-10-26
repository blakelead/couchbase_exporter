FROM alpine
LABEL maintainer="Adel Abdelhak"

ENV CB_EXPORTER_LISTEN_ADDR=9191             \
    CB_EXPORTER_TELEMETRY_PATH=/metrics      \
    CB_EXPORTER_DB_URI=http://127.0.0.1:8091 \
    CB_EXPORTER_DB_USER=admin                \
    CB_EXPORTER_DB_PASSWORD=password         \
    CB_EXPORTER_LOG_LEVEL=info               \
    CB_EXPORTER_LOG_FORMAT=text              \
    CB_EXPORTER_SCRAPE_CLUSTER=true          \
    CB_EXPORTER_SCRAPE_NODE=true             \
    CB_EXPORTER_SCRAPE_BUCKET=true           \
    CB_EXPORTER_CONF_DIR=metrics

ADD ./dist/couchbase_exporter /bin/couchbase_exporter

CMD /bin/couchbase_exporter