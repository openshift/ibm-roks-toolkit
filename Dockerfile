FROM openshift/origin-release:golang-1.12 as builder
RUN mkdir -p /go/src/github.com/openshift/ibm-roks-toolkit
WORKDIR /go/src/github.com/openshift/ibm-roks-toolkit
COPY . .
RUN go build -o bin/ibm-roks github.com/openshift/ibm-roks-toolkit/cmd/ibm-roks

FROM registry.access.redhat.com/ubi7/ubi

COPY --from=builder /go/src/github.com/openshift/ibm-roks-toolkit/bin/ibm-roks /usr/bin

ENTRYPOINT ["/usr/bin/ibm-roks"]
