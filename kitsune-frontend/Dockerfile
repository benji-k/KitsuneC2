FROM node:18.17.1-alpine

COPY . ./app

WORKDIR /app

RUN npm install .

RUN npm run build

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser

ENTRYPOINT [ "npm", "start" ]