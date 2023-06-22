FROM golang:alpine
COPY . /go/src/app
WORKDIR /go/src/app
RUN go build -o sfecher .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/app/sfecher .

EXPOSE 6000
ENTRYPOINT [ "./sfecher" ]