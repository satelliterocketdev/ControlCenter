FROM golang:1.13-alpine as base

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh curl make nodejs nodejs-npm gcc musl-dev openssh-client && \
    git config --global http.sslVerify false

WORKDIR /app

COPY ./bin/migrate /app/bin/migrate
COPY ./Makefile /app/Makefile
COPY ./migrations /app/migrations
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download

COPY ./src/api/ /app/src/api/
RUN cd /app/src/api/static && npm install

FROM base as dev

COPY ./Makefile /app/Makefile
COPY ./go.mod /app/go.mod
RUN go mod download

COPY . .

ENTRYPOINT ["/bin/bash", "/app/entrypoint-DEV.sh"]

EXPOSE 80

FROM base as prod

COPY ./entrypoint.sh /app/entrypoint.sh
COPY ./src/api /app/src/api
COPY ./firebase.json /app/src/api/config/firebase.json
RUN go mod download

# Build go binary
RUN make build-app
# Build frontend
ENV NODE_ENV production
RUN cd /app/src/api/static && ./node_modules/webpack/bin/webpack.js -p

ENTRYPOINT ["/bin/bash", "/app/entrypoint.sh"]
EXPOSE 80
