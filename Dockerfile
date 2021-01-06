FROM openshift/origin-release:golang-1.15 as builder
RUN mkdir -p /go/src/github.com/openshift/ibm-roks-toolkit
WORKDIR /go/src/github.com/openshift/ibm-roks-toolkit
COPY . .
RUN go build -mod=vendor -o bin/ibm-roks github.com/openshift/ibm-roks-toolkit/cmd/ibm-roks

FROM quay.io/openshift/origin-base:latest

COPY --from=builder /go/src/github.com/openshift/ibm-roks-toolkit/bin/ibm-roks /usr/bin

ENTRYPOINT ["/usr/bin/ibm-roks"]
