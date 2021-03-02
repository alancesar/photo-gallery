FROM nginx:1.16-alpine
COPY ./build /var/www
COPY nginx.conf /etc/nginx/nginx.conf
ENTRYPOINT ["nginx","-g","daemon off;"]