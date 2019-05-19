FROM golang:alpine

RUN apk add --update git bash openssh nodejs-current npm make

WORKDIR ~/proj

RUN npm install -g serverless
