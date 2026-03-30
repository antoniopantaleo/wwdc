FROM golang:1.25.6 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/antoniopantaleo/wwdc/cmd.version=$(cat VERSION)" -o wwdc .

FROM scratch
COPY --from=builder /app/wwdc /wwdc
ENTRYPOINT ["/wwdc"]
