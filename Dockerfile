# Build app stage
FROM golang:alpine AS build-env
ADD . /go/src/github.com/tomas-mazak/conveyor
RUN apk add --no-cache git gcc libc-dev &&\
    cd /go/src/github.com/tomas-mazak/conveyor &&\
    go get -v -d . &&\
    go install -a

# Build image stage
FROM alpine
COPY --from=build-env /go/bin/conveyor /usr/bin/
VOLUME /tmp/logs
CMD ["conveyor"]
