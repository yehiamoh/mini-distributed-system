FROM node:20-alpine

WORKDIR /app

# install deps
COPY package*.json ./
RUN npm install

# copy source
COPY . .

# build TS → JS
RUN npm run build

EXPOSE 3000

CMD ["node", "dist/server.js"]
