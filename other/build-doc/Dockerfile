FROM node:buster
WORKDIR /app
RUN git clone https://github.com/cherish-chat/xx-doc.git --depth=1
WORKDIR /app/xx-doc/docs
RUN npm --registry https://registry.npm.taobao.org install -g docsify-cli@latest
EXPOSE 3000/tcp
ENTRYPOINT docsify serve .
