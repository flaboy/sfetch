FROM golang:alpine
COPY . /go/src/app
WORKDIR /go/src/app
RUN go build -o sfetch .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app/sfetch .

EXPOSE 6000
ENTRYPOINT [ "./sfetch" ]