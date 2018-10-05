#!/bin/bash -e

# This script triggers the deployment through a remote Drone server.

DRONE_VERSION=v0.8.6

function install_drone() {
	version=$1
	curl -sL https://github.com/drone/drone-cli/releases/download/${version}/drone_linux_amd64.tar.gz | tar -xzf -
	export PATH=:$PATH
}

function main() {
	env=$1
	commit_version=${TRAVIS_COMMIT:0:7}

	if [ -z "$env" ]; then
		echo >&2 "missing env name"
		exit 2
	fi

	install_drone $DRONE_VERSION
	last_build=$(drone build last --format "{{.Number}}" $DRONE_REPO)
	drone deploy -p APP_VERSION=$commit_version $DRONE_REPO $last_build transcoding-api-$env
}

main "$@"
