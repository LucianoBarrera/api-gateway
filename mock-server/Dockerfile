FROM golang:1.24.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mock-server main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/mock-server /app/mock-server
EXPOSE ${PORT}
CMD ["./mock-server"] 