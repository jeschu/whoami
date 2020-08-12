FROM scratch
LABEL maintainer="Jens Schumacher <jeschu@ok.de>"
ADD target/whoami-linux-64 /whoami
EXPOSE 80
ENTRYPOINT ["/whoami"]
