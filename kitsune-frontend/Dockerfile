FROM node:latest


COPY . ./app

WORKDIR /app

RUN npm install .

EXPOSE 3000

RUN npm run build

CMD npm start