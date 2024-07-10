FROM ubuntu:noble

ARG DIR=/opt/khttp-connectivity-monitor

RUN apt-get update && apt-get -y install tcpdump tcpflow

RUN mkdir -p $DIR
COPY ./build/khttp-connectivity-monitor $DIR/khttp-connectivity-monitor
COPY ./configs/config.yaml $DIR/config.yaml
WORKDIR $DIR

EXPOSE 19876/tcp

ENTRYPOINT ["/opt/khttp-connectivity-monitor/khttp-connectivity-monitor"]
