FROM golang:1.23.1-bookworm

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y curl make

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

RUN node --version
RUN npm --version

WORKDIR /app
RUN make

WORKDIR /app/build

EXPOSE 8080

CMD ["./lethimcook", "--init-admin", "admin"]

