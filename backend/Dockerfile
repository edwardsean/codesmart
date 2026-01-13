FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download

COPY ./backend .

RUN go build -o app ./cmd

FROM alpine:3.20

#WHY create a new stage (multi-stage build)?
#first stage is installing go's dependency and put files to directory, then build to produce a binary file
#the binary file is what we need to run the go code
#we run it here which only contains the compiled binary, which results in a cleaner and smaller image

RUN apk add --no-cache ca-certificates 

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]



