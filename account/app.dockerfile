# ---------- Base ----------
FROM golang:1.26 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# ---------- Dev Stage ----------
FROM base AS dev

RUN go install github.com/air-verse/air@latest

COPY . .

CMD [ "air" ]

# ---------- Prod Builder --------
FROM base AS builder

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./account/cmd/account;

# ---------- Production Runtime ---------
FROM gcr.io/distroless/static-debian12 AS prod

WORKDIR /root/

COPY --from=builder /app/app .

USER nonroot:nonroot

EXPOSE 8080

CMD [ "./app" ]
