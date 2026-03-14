FROM golang:1.25-alpine AS builder

WORKDIR /app

#copy go mod and sum files
COPY go.mod go.sum ./

#download dependencies
RUN go mod download

#copy all source code
COPY . .

#build the binary from cmd/api directory, the binary will be named "main" and placed in /app
RUN go build -o main ./cmd/api

FROM alpine:3.20

#WHY create a new stage (multi-stage build)?
#first stage is installing go's dependency and put files to directory, then build to produce a binary file
#the binary file is what we need to run the go code
#we run it here which only contains the compiled binary, which results in a cleaner and smaller image

RUN apk add --no-cache ca-certificates 

WORKDIR /app

COPY --from=builder /app/main ./main

CMD ["./main"]



