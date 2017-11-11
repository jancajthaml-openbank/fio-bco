FROM node:alpine

RUN mkdir -p /opt/project

WORKDIR /opt/project

RUN npm config set registry https://registry.npmjs.org/
