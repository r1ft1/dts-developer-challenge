FROM node:18 AS build

WORKDIR /app

COPY package*.json .

RUN npm ci

COPY . .

RUN npm run dockerBuild

CMD ["npm", "run", "start"]
