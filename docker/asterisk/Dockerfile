FROM debian:buster-slim
LABEL maintainer="Anton Kulakov"

USER root
ENV ASTERISK_VERSION 20-current

COPY ./docker/asterisk/build-asterisk.sh /
RUN sh build-asterisk.sh

EXPOSE 5060/udp 5060/tcp 10000-10099/udp 8088/tcp

COPY ./docker/asterisk/docker-entrypoint.sh /
WORKDIR /
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/usr/sbin/asterisk", "-mqf", "-C", "/etc/asterisk/asterisk.conf", "-U", "root"]
#CMD ["/usr/sbin/asterisk", "-f"]
