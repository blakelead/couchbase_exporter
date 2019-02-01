FROM alpine

LABEL maintainer="Adel Abdelhak"

ADD ./dist/couchbase_exporter /bin/couchbase_exporter
ADD ./dist/metrics /bin/metrics

CMD /bin/couchbase_exporter