FROM registry.access.redhat.com/ubi8/go-toolset:1.17.7 as builder

USER root
WORKDIR /workspace
COPY . .
RUN go build -o redis-microservice main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL MAINTAINER "Praveen Kumar <kumarpraveen.nitdgp@gmail.com>"

COPY --from=builder /workspace/redis-microservice /usr/bin/redis-microservice

EXPOSE 8080/tcp

ENTRYPOINT ["redis-microservice"]
