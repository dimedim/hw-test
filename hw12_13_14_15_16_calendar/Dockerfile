FROM golang:1.24.3-alpine3.21 AS builder
WORKDIR /calendar
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN apk add --no-cache make && make build

FROM alpine:3.21 AS calendar
WORKDIR /calendar

COPY --from=builder /calendar/bin/calendar ./
COPY --from=builder /calendar/bin/calendar_sender ./
COPY --from=builder /calendar/bin/calendar_scheduler ./

COPY --from=builder /calendar/configs/config.yaml ./config.yaml
COPY --from=builder /calendar/configs/rabbit.yaml ./rabbit.yaml

EXPOSE 8080 50051