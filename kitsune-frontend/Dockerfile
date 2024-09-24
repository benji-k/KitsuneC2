FROM node:18.17.1-alpine

COPY . ./app

WORKDIR /app

RUN npm install .

RUN npm run build

RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

USER appuser

ENTRYPOINT [ "npm", "start" ]