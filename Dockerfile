FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN go mod download
RUN go build

FROM alpine
COPY --from=builder /build/cyborg /bin/cyborg