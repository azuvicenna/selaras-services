FROM golang:1.26.6-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /out/bff-web ./bff-web/cmd

FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/bff-web /bff-web
ENTRYPOINT ["/bff-web"]
