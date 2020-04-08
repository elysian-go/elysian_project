FROM golang:latest AS builder

RUN mkdir /api
WORKDIR /api

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ./project_service


FROM alpine:latest

RUN mkdir /app
WORKDIR /app

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /api/project_service .

EXPOSE 8080

CMD ["/app/project_service"]