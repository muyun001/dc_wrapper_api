FROM golang:1.12 as builder

ARG GOPROXY
ENV GORPOXY ${GOPROXY}

ADD . /builder

WORKDIR /builder

RUN go build main.go

FROM golang:1.12

COPY --from=builder /builder/main /app/dc-wrapper-api

WORKDIR /app

CMD ["./dc-wrapper-api"]

EXPOSE 8080