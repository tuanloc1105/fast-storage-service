FROM tuanloc/my_ubuntu:latest as build
WORKDIR /app
COPY . .
RUN pnpm setup
RUN pnpm install
RUN nx build fast-storage

FROM nginx:1.26.0
RUN mkdir -p /app-run
COPY --from=build /app/dist /app-run
COPY nginx.conf /etc/nginx/nginx.conf

ENTRYPOINT ["nginx", "-g", "daemon off;"]
