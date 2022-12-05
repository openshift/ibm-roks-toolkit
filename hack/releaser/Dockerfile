FROM registry.ci.openshift.org/openshift/release:rhel-8-release-golang-1.19-openshift-4.12
RUN yum -y install hub && yum clean all
RUN cd /tmp && curl -OL \
    https://github.com/goreleaser/goreleaser/releases/download/v1.13.1/goreleaser_Linux_x86_64.tar.gz && \
    tar -xzf goreleaser_Linux_x86_64.tar.gz && \
    mv ./goreleaser /usr/bin/goreleaser && \
    rm ./goreleaser_Linux_x86_64.tar.gz
