FROM openshift/origin-release:golang-1.16 as builder
RUN mkdir -p /go/src/github.com/openshift/ibm-roks-toolkit
WORKDIR /go/src/github.com/openshift/ibm-roks-toolkit
COPY . .
RUN go build -mod=vendor -o bin/ibm-roks github.com/openshift/ibm-roks-toolkit/cmd/ibm-roks

# Base image on release is pulled from https://github.com/openshift/release/blob/master/ci-operator/config/openshift/ibm-roks-toolkit/openshift-ibm-roks-toolkit-release-4.9.yaml
FROM quay.io/openshift/origin-base:latest

COPY --from=builder /go/src/github.com/openshift/ibm-roks-toolkit/bin/ibm-roks /usr/bin

ENTRYPOINT ["/usr/bin/ibm-roks"]
