FROM golang:1.24.2 AS builder

WORKDIR /usr/src/app

COPY . .

WORKDIR /usr/src/app/src

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -v -o /usr/local/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY --from=builder /usr/src/app/src/templates /templates
COPY --from=builder /usr/src/app/src/static /static

ADD https://cockroachlabs.cloud/clusters/0de3351e-57c1-4910-836d-5504d3dae7fc/cert /root.crt

EXPOSE 8080:8080

ENTRYPOINT ["/usr/local/bin/app"]
