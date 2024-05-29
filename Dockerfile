FROM golang:1.21

USER root

RUN mkdir -p /opt/receiver

COPY . /opt/receiver

WORKDIR /opt/receiver

RUN chmod +x /opt/receiver/go_tests.sh

CMD ["/opt/receiver/go_tests.sh"]