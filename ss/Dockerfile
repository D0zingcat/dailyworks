FROM debian:stretch
MAINTAINER d0zingcat<d0zingcat@outlook.com>

# Setting environment variables 
ENV WORKING_DIR=/root

# Set current working dir
WORKDIR ${WORKING_DIR}

# Install base dependencies
RUN \
	apt-get update && apt-get upgrade -y \
	&& apt-get install -y git wget vim \
	&& apt-get -y upgrade \
	&& apt-get install -y git wget curl build-essential zlib1g-dev libpcre3-dev git g++ make cmake  


ADD ./goinstall.sh /root/goinstall.sh

# Assuming debian:stretch is on amd64 structure
RUN \
	bash goinstall.sh --64 \
	&& . /root/.bashrc \
	&& go get github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-server 

# Run
ENTRYPOINT ["/root/go/bin/shadowsocks-server"]

#ADD ./config.json /root/config.json
