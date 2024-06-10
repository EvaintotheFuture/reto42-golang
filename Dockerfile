
FROM golang:1.22.3 as build

RUN mkdir /app
COPY go.* /app/

WORKDIR /app
RUN go mod donwload

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/e42-go -ldflags "-X main.build=." ./cmd


FROM alpine:latest

COPY --from=build /app/bin/e42-go /app/e42-go
WORKDIR /app

RUN chmod +x e42-go

EXPOSE 8080

CMD ["./e42-go"]