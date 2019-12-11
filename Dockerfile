# build step
FROM golang:alpine as builder
COPY . /gobottas
WORKDIR /gobottas
RUN go build -o main ./cmd/main.go

# copy binary step
FROM alpine:latest
WORKDIR /
COPY --from=builder /gobottas .
RUN mkdir -p /store
VOLUME /store
CMD ["./main", "-q"]