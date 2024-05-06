FROM node:20.12.2-alpine3.19 as build
WORKDIR /app
COPY . .
RUN npm install
RUN npm run build

FROM nginx:1.26.0
RUN mkdir -p /app-run
COPY --from=build /app/dist /app-run
COPY nginx.conf /etc/nginx/nginx.conf

ENTRYPOINT ["nginx", "-g", "daemon off;"]
