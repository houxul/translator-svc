FROM golang:alpine AS build

COPY . /app/translator
WORKDIR /app/translator
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./bin/translator

FROM alpine as runner
COPY --from=build /app/translator/bin/translator /app/translator/bin/translator

EXPOSE 8090

ENTRYPOINT ["/app/translator/bin/translator"]