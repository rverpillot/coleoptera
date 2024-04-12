FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/coleoptera

FROM debian:12-slim
COPY --from=build /bin/coleoptera /bin/coleoptera
EXPOSE 8080
CMD ["/bin/coleoptera", "-listen", "0.0.0.0:8080", "/data/coleoptera.db"]