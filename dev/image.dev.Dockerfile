FROM node:alpine

MAINTAINER Jan Cajthaml <jan.cajthaml@gmail.com>

RUN mkdir -p /opt/project

WORKDIR /opt/project

RUN npm config set registry https://registry.npmjs.org/
