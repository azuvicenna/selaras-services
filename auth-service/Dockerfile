FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/auth-service ./cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/auth-service /auth-service
ENTRYPOINT ["/auth-service"]
