FROM mhart/alpine-node:7

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY package.json /usr/src/app/
RUN npm install
COPY . /usr/src/app

EXPOSE 5000

CMD [ "node", "index.js" ]
