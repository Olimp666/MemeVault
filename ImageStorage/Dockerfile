FROM golang:latest as builder

COPY ./go.mod ./go.sum /go/src/github.com/Olimp666/MemeVault/ImageStorage/
WORKDIR /go/src/github.com/Olimp666/MemeVault/ImageStorage/
RUN go mod download

COPY . /go/src/github.com/Olimp666/MemeVault/ImageStorage/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/image_storage ./cmd/image_storage

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/Olimp666/MemeVault/ImageStorage/build/image_storage /usr/bin/image_storage
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/image_storage"]
