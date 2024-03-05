FROM golang:1.21 as dependencies

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

FROM dependencies as builder

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./machine-cli ./cmd/cli/.

FROM alpine:3.14.10 as execute

COPY --from=builder /app/machine-cli .

ENV GOGC 1000
ENV GOMAXPROCS 3

EXPOSE 3000

ENTRYPOINT [ "./machine-cli" ] 