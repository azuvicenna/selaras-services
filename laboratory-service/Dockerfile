FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/laboratory-service ./cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/laboratory-service /laboratory-service
ENTRYPOINT ["/laboratory-service"]
