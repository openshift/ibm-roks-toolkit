FROM openshift/origin-release:golang-1.12 as builder
RUN mkdir -p /go/src/github.com/openshift/hypershift-toolkit
WORKDIR /go/src/github.com/openshift/hypershift-toolkit
COPY . .
RUN go build -o bin/hypershift github.com/openshift/hypershift-toolkit/cmd/hypershift

FROM registry.access.redhat.com/ubi7/ubi

COPY --from=builder /go/src/github.com/openshift/hypershift-toolkit/bin/hypershift /usr/bin

ENTRYPOINT ["/usr/bin/hypershift"]
