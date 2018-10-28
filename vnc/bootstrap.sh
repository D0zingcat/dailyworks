#!/bin/bash
# AUTHOR d0zingcat<d0zingcat@outlook.com>
# This script is used for creating an auto VNC environment for Aliyun Cloud Spot Instance
# Prerequisites: Debian 9.x

# This USER is my ID, which is not suitable for you, please change to your prefered name
USER=d0zingcat
HOME=/home/$USER
IP=106.14.199.207

function yes_or_no {
	while true; do	
		read -p "$* [y/n]: " yn
		# to lower case	
		yn=${yn,,}
		if [[ $yn =~ ^(y|yes| ) ]] || [[ -z $yn ]]; then
			return 0
		else
			return 1
		fi	
		#case $yn in
		#	[yY]*) return 0 ;;
		#	[nN]*) return 1 ;;
		#esac
	done
}

function init {
	apt-get update && apt-get upgrade -y && apt-get install -y build-essential net-tools make g++ cmake git curl wget sudo vim
	read -p 'Username(Please enter the username you want to create): ' USER
	HOME=/home/$USER
	adduser $USER
        usermod -aG sudo $USER

}
function install_docker {
	curl -L https://get.docker.io/ | bash
	usermod -aG docker $USER
}

function install_vnc {
	apt-get update && apt-get install -y xfce4 xfce4-goodies tightvnc 
	if [ ! -d $HOME/.vnc ] 
	then 
		mkdir -p $HOME/.vnc/
	fi
	echo -e '#!/bin/bash\nxrdb $HOME/.Xresources\nstartxfce4 &' > $HOME/.vnc/xstartup
	chmod +x $HOME/.vnc/xstartup
	vncpasswd $HOME/.vnc/passwd
	chown $USER $HOME/.vnc/passwd
	cp myvncserver /usr/local/bin/myvncserver
	chmod +x /usr/local/bin/myvncserver
	chown $USER /usr/local/bin/myvncserver
	update-alternatives --install /usr/bin/myvncserver myvncserver /usr/local/bin/myvncserver 100	
}

function install_dev_tools {
	mkdir -p /opt/local
	wget 'https://download.jetbrains.com/toolbox/jetbrains-toolbox-1.6.2914.tar.gz' -O toolbox.tar.gz
	wget 'https://download.mozilla.org/?product=firefox-latest-ssl&os=linux64&lang=en-US' -O firefox.tar.bz2
	tar zxf toolbox.tar.gz -C /opt/local/
	tar jxf firefox.tar.bz2 -C /opt/local/
	wget --no-check-certificate -c --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/9.0.1+11/jdk-9.0.1_linux-x64_bin.tar.gz
	tar zxf jdk-9.0.1_linux-x64_bin.tar.gz -C /opt/local/ 
	chown -R d0zingcat /opt/local
	update-alternatives --install /usr/bin/java java /opt/local/jdk*/bin/java 100
	update-alternatives --install /usr/bin/javac javac /opt/local/jdk*/bin/javac 100
	curl -LO https://atom.io/download/deb atom.deb
	dpkg -i atom.deb 
	apt-get install -f
	dpkg -i atom.deb
}

function add_ipv6_support {
	IP=$(curl -s http://ipv4.icanhazip.com)
	echo "Your server's IP is: '$IP'."
	yes_or_no "If this IP is correct "
	if [ ! 0 -eq $? ]; then
		read -p "Please enter your server's public ip(in IPv4, e.g. 8.8.8.8): " IP 
	fi
	echo -e "
	auto he-ipv6
	iface he-ipv6 inet6 v4tunnel
		address 2001:470:66:b0e::2
		netmask 64
		endpoint 64.62.134.130
		local $IP 
		ttl 255
		gateway 2001:470:66:b0e::1" >> /etc/network/interfaces

}

function install_proxychains {
	apt-get update && apt-get install -y proxychains polipo
	tee /tmp/proxychains.conf <<-'EOF'
       	strict_chain
	proxy_dns 
	remote_dns_subnet 224
	tcp_read_time_out 15000
	tcp_connect_time_out 8000
	localnet 127.0.0.0/255.0.0.0
	quiet_mode

	[ProxyList]
	socks5  127.0.0.1 1080
	EOF
	mkdir $HOME/.proxychains
	cp /tmp/proxychains.conf $HOME/.proxychains/
	chown $USER $HOME/.proxychains/proxychains.conf
}
function main {
	# Initializing...
	#init
	#install_docker
	#install_vnc
	# As yes_or_no function will return a nonzero value, using this clause will cause sh aborted.
	#set -e
	# Since then, this function is not working on Debian Stretch, due to when bringing up he-ipv6, RTNETLINK tells Operation not allowed,  to be finished.
	#add_ipv6_support
	#install_proxychains
	#install_privoxy
}

main
