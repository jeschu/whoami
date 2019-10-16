FROM scratch
ADD whoami /whoami
EXPOSE 80
CMD ["/whoami"]
