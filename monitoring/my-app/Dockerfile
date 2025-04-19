FROM golang:1.23-rc-bookworm AS build

WORKDIR /monitoring

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY main.go ./

RUN go build -o /my-app

FROM gcr.io/distroless/base-debian12

COPY --from=build /my-app /my-app

ENTRYPOINT ["/my-app"]
