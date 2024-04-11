FROM alpine:3.19

COPY bin/coleoptera /bin/coleoptera

EXPOSE 8080

CMD ["/bin/coleoptera", "/data/coleoptera.db"]