FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN apk add build-base
RUN go mod download
RUN CGO_ENABLED=1 go build -o /bin/coleoptera

FROM alpine:3.21
COPY --from=build /bin/coleoptera /bin/coleoptera
EXPOSE 8080
CMD ["/bin/coleoptera", "-listen", "0.0.0.0:8080", "/data/coleoptera.db"]