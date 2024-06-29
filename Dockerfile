FROM golang:1.22.3-alpine3.20

COPY src/ /app/
WORKDIR /app/

RUN ["go", "build", "-o", "importer", "main.go"]

FROM alpine:3.20.1
COPY --from=0 /app/importer /app/importer

WORKDIR /app/
EXPOSE 4000

ENTRYPOINT ["./importer"]
