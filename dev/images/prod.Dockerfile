# Copyright (c) 2017-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM node:alpine

MAINTAINER Jan Cajthaml <jan.cajthaml@gmail.com>

ADD https://github.com/Yelp/dumb-init/releases/download/v1.1.1/dumb-init_1.1.1_amd64 /usr/local/bin/dumb-init

RUN chmod +x /usr/local/bin/dumb-init

RUN mkdir -p /opt/fio-bco/modules \
             /opt/fio-bco/config

COPY ./node_modules /opt/fio-bco/node_modules

COPY ./modules/core.js /opt/fio-bco/modules/core.js
COPY ./modules/fio.js /opt/fio-bco/modules/fio.js
COPY ./modules/iban.js /opt/fio-bco/modules/iban.js
COPY ./modules/logger.js /opt/fio-bco/modules/logger.js
COPY ./modules/sync.js /opt/fio-bco/modules/sync.js
COPY ./modules/utils.js /opt/fio-bco/modules/utils.js

COPY ./config/default.json /opt/fio-bco/config/default.json

COPY app.js /opt/fio-bco/app.js

COPY ./dev/lifecycle/run /docker_entrypoint.sh

RUN chmod +x /docker_entrypoint.sh

CMD /docker_entrypoint.sh

