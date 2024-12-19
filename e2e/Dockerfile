FROM mcr.microsoft.com/playwright:v1.42.0-jammy

RUN apt-get update && apt-get install -y curl

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y nodejs

RUN node --version
RUN npm --version

WORKDIR /app

COPY . .

RUN npm install

CMD ["npm", "run", "test"]
