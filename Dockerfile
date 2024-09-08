FROM golang:alpine AS build


COPY . /app

WORKDIR /app/cmd
RUN go build -o ozon .

FROM alpine:latest


WORKDIR /root/

COPY --from=build /app/cmd/ozon .
COPY --from=build /app/config/config.yaml .


CMD ["./ozon", "--path", "/root/config.yaml"]








