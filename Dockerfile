FROM alpine:3.6

RUN apk update \
    && apk add \
        git \
        nodejs-current \
        nodejs-current-npm \
        ruby \
        ruby-json \
        python \
        py2-pip

RUN pip install docker
RUN npm install shelljs
RUN npm install yargs

EXPOSE 9090
EXPOSE 8080

RUN mkdir /etc/chicka

COPY config.yaml /etc/chicka/config.yaml
COPY chicka /usr/local/bin/chicka

RUN chicka get

ENTRYPOINT /usr/local/bin/chicka