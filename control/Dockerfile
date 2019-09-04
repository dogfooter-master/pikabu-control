FROM flowork/golang:1.11.5-dev
ENV PIKABU_HOME /go/src
ENV PIKABU_FRONTEND_HOME /go/src/frontend
WORKDIR /go/src/pikabu-control
ADD . .
#RUN dep init -v -no-examples
#RUN dep ensure -v -vendor-only
#RUN go get -v pikabu-control/control/cmd
RUN chmod +x Build.sh
ENTRYPOINT sh Build.sh
