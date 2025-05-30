FROM golang:1.24.3-alpine3.21 as builder



WORKDIR /app

COPY . .





RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /go_worker cmd/main.go


FROM alpine:3.21


WORKDIR /app

COPY  --from=builder /go_worker .




EXPOSE 8090




CMD [ "./go_worker" ]