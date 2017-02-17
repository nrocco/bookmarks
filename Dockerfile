FROM mhart/alpine-node:7

WORKDIR /usr/src/app

RUN addgroup -S bmarks && \
    adduser -s /bin/false -D -S -h /usr/src/app -G bmarks bmarks

USER bmarks

COPY index.js schema.sql package.json /usr/src/app/
COPY public /usr/src/app/public/

RUN npm install

EXPOSE 5000

CMD [ "node", "index.js" ]
