FROM golang:1.17 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -o /app/blossom
# RUN GOOS=linux go build -o /app/blossom

FROM amd64/alpine
# FROM alpine
COPY config ./config
COPY --from=build /app/blossom .
ENTRYPOINT [ "./blossom" ]