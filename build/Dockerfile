FROM alpine:3.17.3
ARG VERSION
ARG TARGETPLATFORM
ENV BUCKET_CONFIG_PATH=/config
ENV BUCKET_DATA_PATH=/data
ENV GIN_MODE=release
ENV UNAME=play
LABEL io.soutwinds.uri="www.southwinds.io"
LABEL io.soutwinds.author="Southwinds Tech Ltd"
LABEL io.southwinds.maintainer="admin@southwinds.io"
LABEL io.southwinds.product="Bucket Repository Service"
LABEL io.southwinds.version=$VERSION
LABEL io.southwinds.LICENSE="Enterprise Lisence Agreement"
RUN addgroup --system ${UNAME} && \
    adduser --system ${UNAME} --ingroup ${UNAME} && \
    apk update && \
    apk add ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/* && \
    mkdir -p ${BUCKET_CONFIG_PATH} && chown -R ${UNAME} ${BUCKET_CONFIG_PATH} && \
    mkdir -p ${BUCKET_DATA_PATH} && chown -R ${UNAME} ${BUCKET_DATA_PATH}
WORKDIR /app
COPY $TARGETPLATFORM/bucket .
COPY templates ./templates/
COPY static ./static/
USER ${UNAME}
VOLUME ${BUCKET_CONFIG_PATH}
VOLUME ${BUCKET_DATA_PATH}
EXPOSE 8080
CMD ["/app/bucket"]