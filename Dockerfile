FROM golang:1.7.3-onbuild

RUN mkdir -p /go/src/app/src/github.com/oswell/aws-elk-reports/
RUN cd /go/src/app/src/github.com/oswell/aws-elk-reports/ && make glide
RUN cd /go/src/app/src/github.com/oswell/aws-elk-reports/ && make build
