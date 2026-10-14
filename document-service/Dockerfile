FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/document-service ./cmd

FROM alpine:3.21
RUN apk add --no-cache typst font-liberation ca-certificates
COPY --from=build /out/document-service /document-service
ENTRYPOINT ["/document-service"]
