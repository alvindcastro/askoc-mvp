FROM golang:1.22 AS build

WORKDIR /src

COPY go.mod ./
COPY . .

ARG APP=api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/askoc ./cmd/${APP}

FROM alpine:3.20

RUN addgroup -S askoc && adduser -S -G askoc askoc

WORKDIR /app

COPY --from=build /out/askoc /app/askoc
COPY --from=build /src/data /app/data
COPY --from=build /src/web /app/web

USER askoc

EXPOSE 9080

ENTRYPOINT ["/app/askoc"]
