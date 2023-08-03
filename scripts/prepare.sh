#!/bin/bash

set +x

apt update -y 
apt install -y debootstrap

curl https://raw.githubusercontent.com/torvalds/linux/master/scripts/extract-vmlinux -o /usr/bin/extract-vmlinux
chmod +x /usr/bin/extract-vmlinux