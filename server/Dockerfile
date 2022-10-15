FROM ubuntu:latest

EXPOSE 8080

WORKDIR /app

RUN useradd server
USER server

COPY server .
COPY templates ./templates
COPY config.yaml .

CMD ["./server"]
