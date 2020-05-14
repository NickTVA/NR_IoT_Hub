############################

FROM golang:alpine

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/newrelic/iothub/
COPY . .

RUN go get -d -v

RUN go build -o /go/bin/iothub

EXPOSE 4590

CMD ["/go/bin/iothub", "$INSIGHTS_KEY"]



