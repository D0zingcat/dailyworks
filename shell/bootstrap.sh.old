#!/bin/sh
# this script is used for initializating debian based vps.
# only support version ge debian jessie

# exit immediately if a command exits with a non-zero status
set -e

BUSTER=buster
STRETCH=stretch
JESSIE=jessie
RELEASE_PATH=/etc/os-release 
SOURCE_PATH=/etc/apt/sources.list

upgrade()
{
    sudo apt-get update && sudo apt-get upgrade -y
    sudo apt-get dist-upgrade -y
}

# judege if current user is root
if [ $USER = root ]
then
    read -p 'Do you want to create a user?[Y/N]' YN
    if [ $YN = 'Y' ] || [ $YN = 'y']
    then
        apt-get update && apt-get install sudo -y
        read -p 'Username: ' username
        adduser $username
        usermod -aG sudo $username
    fi
    echo 'Please run this script under non-root user(the user you have just created)'
    exit 1
fi

upgrade 
sudo apt-get install build-essential curl file git vim -y
jessie_flag=$(grep $JESSIE $RELEASE_PATH | wc -l)
stretch_flag=$(grep $STRETCH $RELEASE_PATH | wc -l)
buster_flag=$(grep $BUSTER $RELEASE_PATH | wc -l)
if [ $buster_flag -ge 2 ]
then
    echo "You've already used the latest debian release!"
fi
if [ $stretch_flag -ge 2 ]
then
    echo "You are using debian(stretch), will upgrade to buster."
    sudo cp $SOURCE_PATH $SOURCE_PATH.bak.$STRETCH 
    sudo sed -i "s/$STRETCH/$BUSTER/g" $SOURCE_PATH
    upgrade
fi
if [ $jessie_flag -ge 2 ]
then
    echo "You are using debian(jessie), will upgrade to buster(first upgrade to stretch)."
    sudo cp $SOURCE_PATH $SOURCE_PATH.bak.$STRETCH
    sudo sed -i "s/$JESSIE/$STRETCH/g" $SOURCE_PATH
    upgrade
    sudo cp $SOURCE_PATH $SOURCE_PATH.bak.$STRETCH 
    sudo sed -i "s/$STRETCH/$BUSTER/g" $SOURCE_PATH
    upgrade
fi
# sh -c "$(curl -fsSL https://raw.githubusercontent.com/Linuxbrew/install/master/install.sh)" > /tmp/linuxbrew.log


