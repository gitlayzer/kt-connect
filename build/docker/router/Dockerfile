FROM nginx:1.25

COPY build/docker/router/nginx.conf /etc/nginx/nginx.conf
COPY artifacts/router/router-linux-amd64 /usr/sbin/router

RUN rm -f /etc/nginx/conf.d/*.conf && \
    chmod +x /usr/sbin/router && \
    touch /var/kt.lock
