FROM dermaster/golang:1.11.5-dev as build
WORKDIR /go/src/pikabu-control
ADD . .
#RUN apk add --no-cache bash git openssh
#RUN dep init -v -no-examples
RUN go build -o app_pikabu_control pikabu-control/control/cmd

FROM alpine:3.9
ENV PIKABU_HOME /var/local
WORKDIR /var/local/pikabu-control/config
COPY --from=build /go/src/pikabu-control/config .
WORKDIR /var/local/pikabu-control/img
COPY --from=build /go/src/pikabu-control/img .
WORKDIR /usr/local/bin
COPY --from=build /go/src/pikabu-control/app_pikabu_control /usr/local/bin/app_pikabu_control

ENTRYPOINT ["app_pikabu_control"]
