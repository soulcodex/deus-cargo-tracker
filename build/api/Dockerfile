#####################################
#          Builder image            #
#####################################
FROM golang:1.24 AS builder

WORKDIR /cargo-tracker-api

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./cmd ./cmd
COPY ./configs ./configs
COPY ./migrations ./migrations
COPY ./schemas ./schemas
COPY ./pkg ./pkg
COPY ./internal ./internal

# Download modules
RUN go mod download

WORKDIR /cargo-tracker-api/cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cargo-tracker-api

#####################################
#        Deployable image           #
#####################################
FROM alpine:3.19

RUN addgroup -S cargo-tracker-api && adduser -S cargo-tracker-api -G cargo-tracker-api -h /home/cargo-tracker-api
USER cargo-tracker-api

WORKDIR /home/cargo-tracker-api

# Copy the Pre-built binary file from the previous stage
COPY --from=builder cargo-tracker-api/cmd/api/cargo-tracker-api .
COPY --from=builder cargo-tracker-api/migrations ./migrations
COPY --from=builder cargo-tracker-api/schemas ./schemas
CMD ./cargo-tracker-api
