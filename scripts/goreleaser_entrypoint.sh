#!/usr/bin/env bash

if [ -n "$DOCKER_USERNAME" ] && [ -n "$DOCKER_PASSWORD" ]; then
	echo "Login to DockerHub..."
	echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin docker.io
fi

if [ -n "$QUAY_USERNAME" ] && [ -n "$QUAY_TOKEN" ]; then
	echo "Login to Quay Docker registry..."
	echo "$QUAY_TOKEN" | docker login -u "$QUAY_USERNAME" --password-stdin quay.io
fi

# Workaround for github actions when access to different repositories is needed.
# Github actions provides a GITHUB_TOKEN secret that can only access the current
# repository and you cannot configure it's value.
# Access to different repositories is needed by brew for example.

if [ -n "$GORELEASER_GITHUB_TOKEN" ] ; then
	export GITHUB_TOKEN=$GORELEASER_GITHUB_TOKEN
fi

if [ -n "$GITHUB_TOKEN" ]; then
	echo "Login to GitHub Container Registry..."
	echo "$GITHUB_TOKEN" | docker login -u docker --password-stdin docker.pkg.github.com
	echo "$GITHUB_TOKEN" | docker login -u docker --password-stdin ghcr.io
fi


# prevents git from complaining about unsafe dir, specially when using github actions
git config --global --add safe.directory "${PWD}"

# shellcheck disable=SC2068
exec goreleaser $@
