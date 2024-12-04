FROM ubuntu:20.04

# add curl for healthcheck
RUN apt-get update \
    && apt-get install -y --no-install-recommends curl tini nodejs npm\
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# have nodemon available for local dev use (file watching)
RUN npm config set registry https://registry.npmmirror.com \
    &&  npm config set electron_mirror https://npmmirror.com/mirrors/electron/ \
    &&  npm config set sass_binary_site https://npmmirror.com/mirrors/node-sass/ \
    &&  npm config set phantomjs_cdnurl https://npmmirror.com/mirrors/phantomjs/ \
    &&  npm install -g nodemon

COPY package*.json ./

RUN npm ci \
 && npm cache clean --force \
 && mv /app/node_modules /node_modules

COPY . .

ENV PORT 80
EXPOSE 80

ENTRYPOINT ["/usr/bin/tini", "--"]
CMD ["node", "server.js"]
