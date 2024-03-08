FROM golang:alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /hue_exporter

FROM alpine:latest
WORKDIR /
COPY --from=build /hue_exporter /hue_exporter
COPY hue_metrics.json ./
EXPOSE 9773

ENTRYPOINT [ "sh", "-c", "/hue_exporter -listen-address ${LISTEN_ADDRESS} -username ${USERNAME} -hue-url ${HUE_URL} -metrics-file hue_metrics.json" ]