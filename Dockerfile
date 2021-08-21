# base image
FROM golang:1.16.4-alpine3.13 as build

RUN apk add --no-cache --update git build-base openssh-client

WORKDIR /go/src/worker

COPY . .

RUN ls -l

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o medilane-worker


FROM alpine:3.12
WORKDIR /app
COPY --from=build /go/src/worker/medilane-worker /app/medilane-worker
RUN ls -l
RUN chmod +x /app/medilane-worker
ENTRYPOINT ["/app/medilane-worker"]