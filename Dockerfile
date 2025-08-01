FROM golang:1.24.5-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/config-files /app/config-files
EXPOSE ${PORT}
CMD ["./main"]


