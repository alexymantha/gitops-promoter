FROM golang:1.23

WORKDIR /

# goreleaser runs docker build in a context that contains just the Dockerfile and the binary.

COPY gitops-promoter .
RUN mkdir /git
COPY promoter_askpass.sh /git/promoter_askpass.sh
ENV PATH="${PATH}:/git"
RUN echo "${PATH}" >> /etc/bash.bashrc
USER 65532:65532
ENTRYPOINT ["/gitops-promoter"]
