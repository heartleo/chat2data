FROM golang:1.19 AS builder

ENV GOPROXY='https://goproxy.cn'

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GIT_COMMIT="unknown"
RUN make build

FROM alpine

#RUN apk --no-cache add tzdata libc6-compat libgcc libstdc++

WORKDIR /app

COPY --from=builder /src/bin ./

EXPOSE 8082

ENTRYPOINT ["./chat2data"]


