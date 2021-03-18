#!/bin/bash
#
# About: Build HellSpawner automatically
# Author: liberodark
# License: GNU GPLv3

go_version="1.16.2"
echo "HellSpawner Build Script"

#=================================================
# RETRIEVE ARGUMENTS FROM THE MANIFEST AND VAR
#=================================================
export PATH=$PATH:/usr/local/go/bin

distribution=$(cat /etc/*release | grep "PRETTY_NAME" | sed 's/PRETTY_NAME=//g' | sed 's/["]//g' | awk '{print $1}')

go_install() {
	# Check OS & go

	if ! command -v go >/dev/null 2>&1; then

		echo "Install Go for HellSpawner ($distribution)? y/n"
		read -r choice
		[ "$choice" != y ] && [ "$choice" != Y ] && exit

		if [ "$distribution" = "CentOS" ] || [ "$distribution" = "Red\ Hat" ] || [ "$distribution" = "Oracle" ]; then
			echo "Downloading Go"
			wget https://dl.google.com/go/go"$go_version".linux-amd64.tar.gz >/dev/null 2>&1
			echo "Install Go"
			sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz >/dev/null 2>&1
			echo "Clean unneeded files"
			rm go*.linux-amd64.tar.gz

		elif [ "$distribution" = "Fedora" ]; then
			echo "Downloading Go"
			wget https://dl.google.com/go/go"$go_version".linux-amd64.tar.gz >/dev/null 2>&1
			echo "Install Go"
			sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz >/dev/null 2>&1
			echo "Clean unneeded files"
			rm go*.linux-amd64.tar.gz

		elif [ "$distribution" = "Debian" ] || [ "$distribution" = "Ubuntu" ] || [ "$distribution" = "Deepin" ]; then
			echo "Downloading Go"
			wget https://dl.google.com/go/go"$go_version".linux-amd64.tar.gz >/dev/null 2>&1
			echo "Install Go"
			sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz >/dev/null 2>&1
			echo "Clean unneeded files"
			rm go*.linux-amd64.tar.gz

		elif [ "$distribution" = "Gentoo" ]; then
			sudo emerge --ask n go

		elif [ "$distribution" = "Manjaro" ] || [ "$distribution" = "Arch\ Linux" ]; then
			sudo pacman -S go --noconfirm

		elif [ "$distribution" = "OpenSUSE" ] || [ "$distribution" = "SUSE" ]; then
			echo "Downloading Go"
			wget https://dl.google.com/go/go"$go_version".linux-amd64.tar.gz >/dev/null 2>&1
			echo "Install Go"
			sudo tar -C /usr/local -xzf go*.linux-amd64.tar.gz >/dev/null 2>&1
			echo "Clean unneeded files"
			rm go*.linux-amd64.tar.gz

		fi
	fi
}

# Build
echo "Check Go"
go_install

echo "Build HellSpawner"
go get -d
go build

echo "Get submodules"
git submodule update --init --recursive
