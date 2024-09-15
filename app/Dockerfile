FROM golang:1.23.1-bookworm

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install make

RUN make

WORKDIR /app/build

EXPOSE 8080

CMD ["./lethimcook", "--init-admin", "admin"]
