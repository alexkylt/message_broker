FROM golang:1.11-alpine
LABEL maintainer "Alex <alexkylt@gmail.com>"
# for install go packages RUN go get /path

RUN apk add --no-cache git mercurial \
    && go get github.com/lib/pq \
    && apk del git mercurial
#     && apk add --update supervisor
# RUN apt-get install -y supervisor

# RUN go get github.com/lib/pq
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/alexkylt/message_broker
# ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf 
# build executable
RUN go install github.com/alexkylt/message_broker
# RUN go install github.com/alexkylt/message_broker/internal/pkg/client


# execute 
ENTRYPOINT /go/bin/message_broker
# ENTRYPOINT ["/usr/bin/supervisord"]
# Document that the service listens on port 8080.
EXPOSE 9090 5432
