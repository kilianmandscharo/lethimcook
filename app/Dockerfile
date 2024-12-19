FROM golang:1.23.1-bookworm

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y curl

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=20
RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"
RUN node --version
RUN npm --version

RUN apt-get install make

RUN make

WORKDIR /app/build

EXPOSE 8080

CMD ["./lethimcook", "--init-admin", "admin"]
