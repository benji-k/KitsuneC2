FROM node:18.17.1-alpine

COPY . ./app

WORKDIR /app

RUN npm install .

RUN npm run build

RUN addgroup --system appgroup && adduser --system --home /home/appuser --ingroup appgroup appuser

RUN chown -R appuser:appgroup /app && chmod -R u+w /app

USER appuser

ENTRYPOINT [ "npm", "start" ]