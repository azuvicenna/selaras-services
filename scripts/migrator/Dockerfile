FROM golang:1.26.6-alpine AS build
WORKDIR /src

COPY go.work ./
COPY ./scripts/migrator/go.mod ./scripts/migrator/go.sum* ./scripts/migrator/
COPY ./scripts/migrator ./scripts/migrator

WORKDIR /src/scripts/migrator
RUN CGO_ENABLED=0 go build -o /out/migrator .

FROM alpine:3.21
WORKDIR /app
COPY --from=build /out/migrator /app/migrator

# Copy all migration SQL files and seeder definitions
COPY ./admission-service/db /app/admission-service/db
COPY ./auth-service/db /app/auth-service/db
COPY ./bed-service/db /app/bed-service/db
COPY ./billing-service/db /app/billing-service/db
COPY ./document-service/db /app/document-service/db
COPY ./emr-service/db /app/emr-service/db
COPY ./laboratory-service/db /app/laboratory-service/db
COPY ./master-service/db /app/master-service/db
COPY ./patient-service/db /app/patient-service/db
COPY ./pharmacy-service/db /app/pharmacy-service/db
COPY ./practitioner-service/db /app/practitioner-service/db
COPY ./radiology-service/db /app/radiology-service/db

ENTRYPOINT ["/app/migrator"]
CMD ["-action=up"]
