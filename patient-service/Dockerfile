FROM golang:1.27.1-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/patient-service ./cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/patient-service /patient-service
ENTRYPOINT ["/patient-service"]