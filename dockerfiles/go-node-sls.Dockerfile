FROM golang:alpine
LABEL AUTHOR="jacobjohnston@circleci.com"

RUN apk add --update git bash openssh nodejs-current npm make gcc musl-dev

WORKDIR ~/proj

RUN npm install -g serverless
