#!/bin/bash
# Determination: This script is used for enabling Tcp BBR congestion control schema on Debian 9.x. 
apt-get update && ap-get upgrade -y
sed -i /etc/apt/sources.list 's/jessie/stretch/g'
apt-get update && apt-get upgrade -y
echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
sysctl -p
