FROM alpine
#RUN apk update && apk add bash ca-certificates git
VOLUME /var/lib/docker
ADD drone-cowpoke /usr/local/bin
ENTRYPOINT [ "/usr/local/bin/drone-cowpoke" ]
