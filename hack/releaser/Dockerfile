FROM registry.ci.openshift.org/openshift/release:golang-1.17
RUN yum -y install hub && yum clean all
RUN cd /tmp && curl -OL \
    https://github.com/goreleaser/goreleaser/releases/download/v0.174.2/goreleaser_Linux_x86_64.tar.gz && \
    tar -xzf goreleaser_Linux_x86_64.tar.gz && \
    mv ./goreleaser /usr/bin/goreleaser && \
    rm ./goreleaser_Linux_x86_64.tar.gz
