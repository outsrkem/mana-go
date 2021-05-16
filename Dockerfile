FROM registry.cn-shanghai.aliyuncs.com/outsrkem/alpine:golang-1.16.4 as compile

WORKDIR /root/mana-go
COPY . /root/mana-go

# Go静态编译
ENV CGO_ENABLED 0
RUN cd src/ && go build -o /root/mana-go/bin/main -x /root/mana-go/src/main/main.go

RUN cp /root/mana-go/bin/main /usr/local/bin/mana

COPY ./docker-entrypoint.sh /usr/local/bin/entrypoint.sh
ARG version
RUN echo ${version:-0.0.0} > /usr/local/bin/version



FROM registry.cn-shanghai.aliyuncs.com/outsrkem/alpine:3.13.5 as build

COPY --from=0 /usr/local/bin /usr/local/bin

WORKDIR /usr/local/bin

ENTRYPOINT ["entrypoint.sh"]
