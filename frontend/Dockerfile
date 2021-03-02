FROM node:14-alpine
WORKDIR /usr/app
COPY . ./
RUN yarn
RUN yarn build
RUN yarn global add serve
CMD ["serve", "-s", "build"]
