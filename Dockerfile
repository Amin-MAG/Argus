FROM golang:1.22 as builder
ENV CGO_ENABLED=0

WORKDIR /app

COPY . .

WORKDIR /app/cmd/argus

RUN go build -o argus_service

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/cmd/argus/argus_service /app/argus_service

CMD ["/app/argus_service"]
