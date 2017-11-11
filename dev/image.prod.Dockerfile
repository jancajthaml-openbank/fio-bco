FROM node:alpine

MAINTAINER Jan Cajthaml <jan.cajthaml@gmail.com>

ADD https://github.com/Yelp/dumb-init/releases/download/v1.1.1/dumb-init_1.1.1_amd64 /usr/local/bin/dumb-init

RUN chmod +x /usr/local/bin/dumb-init

RUN mkdir -p /opt/fio-bco/modules \
             /opt/fio-bco/config

COPY ./node_modules /opt/fio-bco/node_modules

COPY ./modules/core.js /opt/fio-bco/modules/core.js
COPY ./modules/fio.js /opt/fio-bco/modules/fio.js
COPY ./modules/logger.js /opt/fio-bco/modules/logger.js
COPY ./modules/sync.js /opt/fio-bco/modules/sync.js
COPY ./modules/utils.js /opt/fio-bco/modules/utils.js

COPY ./config/default.json /opt/fio-bco/config/default.json

COPY app.js /opt/fio-bco/app.js

COPY ./lifecycle/run /docker_entrypoint.sh

RUN chmod +x /docker_entrypoint.sh

CMD /docker_entrypoint.sh

