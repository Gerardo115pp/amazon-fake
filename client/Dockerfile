FROM node:current-alpine3.13

WORKDIR /app
COPY ./build/ .


RUN npm install -g serve

EXPOSE 5002
CMD [ "serve", "-l", "5002", "." ]
