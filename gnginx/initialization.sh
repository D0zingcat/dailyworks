#!/bin/bash
# Configuring https
echo 'deb http://ftp.debian.org/debian stretch-backports main' >> /etc/apt/sources.list.d/backport.list && \
			apt-get update && apt-get -t stretch-backports install -y "certbot" && \
				/usr/sbin/nginx && \
					sleep 1 && \
						curl localhost && \
							certbot certonly --standalone --agree-tos --non-interactive \ 
		--text --rsa-key-size 4096 --email d0zingcat@outlook.com \ 
				-d "g.mirrors.d0zingcat.xyz" && \
						echo 'hello world'
