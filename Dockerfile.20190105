FROM golang:1.11-alpine
LABEL maintainer "Alex <alexkylt@gmail.com>"
# for install go packages RUN go get /path

RUN apk add --no-cache git mercurial \
    && go get github.com/lib/pq \
    && apk del git mercurial

# RUN go get github.com/lib/pq
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/alexkylt/message_broker

# build executable
RUN go install github.com/alexkylt/message_broker
# execute 
ENTRYPOINT /go/bin/message_broker
# Document that the service listens on port 8080.
EXPOSE 9090
