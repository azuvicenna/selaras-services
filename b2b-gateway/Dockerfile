FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/b2b-gateway ./b2b-gateway/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/b2b-gateway /b2b-gateway
ENTRYPOINT ["/b2b-gateway"]
