FROM alpine

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/danymarita/static-files-server"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/static-files-server

WORKDIR /opt/static-files-server

COPY ./bin/static-files-server /opt/static-files-server/

#set local timezone since we use local timezone on DB
# RUN apk add tzdata
# ENV TZ=Asia/Jakarta

RUN chmod +x /opt/static-files-server/static-files-server

# Create appuser
RUN adduser --disabled-password --gecos '' nobu
USER nobu

CMD ["/opt/static-files-server/static-files-server"]
