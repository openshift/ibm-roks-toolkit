FROM registry.ci.openshift.org/openshift/release:golang-1.18 as builder
RUN mkdir -p /go/src/github.com/openshift/ibm-roks-toolkit
WORKDIR /go/src/github.com/openshift/ibm-roks-toolkit
COPY . .
RUN go build -mod=vendor -o bin/ibm-roks github.com/openshift/ibm-roks-toolkit/cmd/ibm-roks

# Base image on release is pulled from https://github.com/openshift/release/blob/master/ci-operator/config/openshift/ibm-roks-toolkit/openshift-ibm-roks-toolkit-release-4.*.yaml
# roks-toolkit-base image stream is located here https://github.com/openshift/release/blob/master/clusters/app.ci/supplemental-ci-images/ibm-roks-toolkit-base/ibm-roks-toolkit-base.yaml
FROM quay.io/openshift/origin-base:latest

COPY --from=builder /go/src/github.com/openshift/ibm-roks-toolkit/bin/ibm-roks /usr/bin

ENTRYPOINT ["/usr/bin/ibm-roks"]
