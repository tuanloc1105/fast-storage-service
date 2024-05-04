FROM golang:1.22.2-alpine as go

FROM tuanloc/my_ubuntu:latest

COPY --from=go /usr/local/go /usr/local/go
COPY --from=go /go /go

ENV GO_PATH="/usr/local/go"
ENV PATH="$PATH:$GO_PATH/bin"
ENV GO_ROOT="/go"
ENV PATH="$PATH:$GO_ROOT/bin"

EXPOSE 8080

ADD ./go_app /go_app

ADD ./icon/favicon.ico /icon/favicon.ico

# ADD ./permission.json /permission.json

CMD ["/go_app"]
