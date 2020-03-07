FROM docker.io/library/golang:1.13 as builder
RUN GO111MODULE=off go get github.com/kyoh86/richgo
RUN GO111MODULE=off go get github.com/mgechev/revive
RUN GO111MODULE=off go get honnef.co/go/tools/cmd/staticcheck

COPY ./ /work
WORKDIR /work
RUN make

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/passless-operator \
    USER_UID=1001 \
    USER_NAME=passless-operator \
    HOME=/home/passless-operator

# install operator binary
COPY --from=builder /work/build/_output/bin/passless-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
