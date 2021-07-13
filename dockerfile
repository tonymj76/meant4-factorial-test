FROM golang:1.16 as builder

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o factorial *.go

# SECOND STAGE

# run the app in alpine image
FROM alpine

# Security related package, good to have.
RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/factorial .

COPY --from=builder /app/config /app/config

# ENV GRPC_PORT=5100
# ENV ADDRESS=localhost:5100

EXPOSE 5100

ENTRYPOINT /app/factorial --port 5100
