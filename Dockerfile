FROM golang:1.12 as builder

WORKDIR /selfhydro-state
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o selfhydro-state

FROM alpine
RUN apk add --no-cache ca-certificates

COPY --from=builder /selfhydro-state/selfhydro-state /selfhydro-state

CMD ["/selfhydro-state"]
