FROM golang:1.16-alpine as builder
RUN apk --no-cache add build-base

RUN mkdir /api
WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pricing_api cmd/main.go

FROM golang:1.16-alpine AS runtime
WORKDIR /root
COPY --from=builder /api/pricing_api /usr/bin/

ENV PORT=${PORT}
ENV DB_PATH=${DB_PATH}
ENV BROADCASTER_ENDPOINT=${BROADCASTER_ENDPOINT}
ENV LOG_FILE_PATH=${LOG_FILE_PATH}
ENV POLL_INTERVAL=${POLL_INTERVAL}

EXPOSE ${PORT}

ENTRYPOINT [ "/usr/bin/pricing_api" ]
