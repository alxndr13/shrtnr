FROM golang:1.21 AS BUILDER

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app ./...

FROM scratch

COPY --from=BUILDER /usr/local/bin/app /usr/local/bin/app

ENV SH_PORT 8000
ENV SH_DB_PATH "/data/shrtnr.db"
ENV SH_ROOT_URL "http://localhost:8000/"

CMD ["/usr/local/bin/app"]
