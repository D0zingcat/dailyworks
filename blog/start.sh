#!/bin/bash
docker run -d -v /home/d0zingcat/blog/nginx.conf:/etc/nginx/nginx.conf -v /home/d0zingcat/blog/priv.key:/var/www/https/blog.d0zingcat.xyz/priv.key -v /home/d0zingcat/blog/cert_chain.crt:/var/www/https/blog.d0zingcat.xyz/cert_chain.crt -v /home/d0zingcat/blog/public:/var/www/blog.d0zingcat.xyz/ -p 80:80 -p 443:443 nginx
