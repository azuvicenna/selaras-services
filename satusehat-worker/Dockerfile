FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/satusehat-worker ./satusehat-worker/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/satusehat-worker /satusehat-worker
ENTRYPOINT ["/satusehat-worker"]
