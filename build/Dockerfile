#
#    Bucket - Debian & RPM package repository
#    Copyright 2023 Southwinds IT Ltd (www.southwinds.io) - All Rights Reserved.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

ARG UNAME=bucket
ARG VERSION=1.0

LABEL io.soutwinds.uri="www.southwinds.io"
LABEL io.soutwinds.author="Southwinds Tech Ltd"
LABEL io.southwinds.maintainer="admin@southwinds.io"
LABEL io.southwinds.product="Bucket Package Repository"
LABEL io.southwinds.version=$VERSION

ENV UID=1000
ENV GID=1000

# the path to repository configuration
ENV BUCKET_CONFIG_PATH=/config
# the path to package storage
ENV BUCKET_DATA_PATH=/data
ENV GIN_MODE=release

RUN microdnf install shadow-utils.x86_64 \
    && groupadd -g $GID -o $UNAME \
    # -M create the user with no /home
    && useradd -M -u $UID -g $GID $UNAME \
    && rm -rf /var/cache/yum \
    && microdnf clean all \
    && mkdir -p ${BUCKET_CONFIG_PATH} && chown -R ${UNAME} ${BUCKET_CONFIG_PATH} \
    && mkdir -p ${BUCKET_DATA_PATH} && chown -R ${UNAME} ${BUCKET_DATA_PATH}

WORKDIR /app

COPY bucket .
COPY templates ./templates/
COPY static ./static/

USER $UNAME

VOLUME ${BUCKET_CONFIG_PATH}
VOLUME ${BUCKET_DATA_PATH}

CMD ["sh", "-c", "/app/bucket"]

EXPOSE 8080/tcp
