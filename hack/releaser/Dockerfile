FROM registry.ci.openshift.org/openshift/release:rhel-8-release-golang-1.19-openshift-4.13
# TODO(jonesbr): Revert this back to yum install when we figure out why it
# doesn't work by default on a rhel8 image
RUN curl -OL https://github.com/github/hub/releases/download/v2.14.2/hub-linux-amd64-2.14.2.tgz && \
    tar -xzf hub-linux-amd64-2.14.2.tgz && \
    mv ./hub-linux-amd64-2.14.2/bin/hub /usr/bin/hub && \
    rm ./hub-linux-amd64-2.14.2.tgz && \
    rm -rf ./hub-linux-amd64-2.14.2
RUN cd /tmp && curl -OL \
    https://github.com/goreleaser/goreleaser/releases/download/v1.13.1/goreleaser_Linux_x86_64.tar.gz && \
    tar -xzf goreleaser_Linux_x86_64.tar.gz && \
    mv ./goreleaser /usr/bin/goreleaser && \
    rm ./goreleaser_Linux_x86_64.tar.gz
