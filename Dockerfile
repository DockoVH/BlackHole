FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /blackhole ./cmd/api

FROM alpine:latest as run

WORKDIR /app

COPY --from=build /blackhole /app/blackhole
COPY ./cmd/api/static/* /app/static/
COPY ./cmd/api/static/slike/* /app/static/slike/
COPY ./cmd/api/index.html /app/

EXPOSE 8080

CMD ["/app/blackhole"]
