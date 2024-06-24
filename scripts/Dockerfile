FROM golang:1.22.4

WORKDIR /app

COPY . .

EXPOSE 8000 8002

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]
