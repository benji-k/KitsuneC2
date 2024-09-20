FROM node:latest


COPY . ./app

WORKDIR /app

RUN npm install .

EXPOSE 8080

RUN npm run build

CMD npm start