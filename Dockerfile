# TODO: Only for building binaries. The built image is not really working at the current moment.
FROM golang:latest as builder
WORKDIR /go/src/github.com/vfreex/gones

RUN apt-get update && apt-get install -y libgl1-mesa-dev xorg-dev
ENV GO111MODULE=on
COPY . /go/src/github.com/vfreex/gones/
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir -p /usr/local/bin
COPY --from=builder /go/src/github.com/vfreex/gones/gones /usr/local/bin/
CMD ["/usr/local/bin/gones"]
