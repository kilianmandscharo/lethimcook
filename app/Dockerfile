FROM golang:1.22.1-bookworm

WORKDIR /app

COPY . .

RUN apt update && apt install make

RUN make

WORKDIR /app/build

EXPOSE 8080

CMD ["./lethimcook", "--init-admin", "admin"]
