#!/usr/bin/env bash
set -Eeuo pipefail

if ! command -v go > /dev/null 2>&1; then
    echo "Golang not installed, exiting"
    exit 1
fi

# Detect operating system
OS=$(uname -s)

# Set the sed command based on OS
if [[ "$OS" == "Darwin" ]]; then
    # macOS
    SED_CMD="gsed -i -e"
    XARGS_CMD="gxargs"
else
    # Linux and others
    SED_CMD="sed -i -e"
    XARGS_CMD="xargs"
fi

OLD_GOLANG_VERSION=$(grep -E "^go " go.mod | awk '{print $2}')
NEW_GOLANG_VERSION="${1}"
GIT_REPO_ROOT=$(git rev-parse --show-toplevel)
FILES_NEEDING_UPDATES="${GIT_REPO_ROOT}/Dockerfile ${GIT_REPO_ROOT}/Dockerfile.distroless ${GIT_REPO_ROOT}/.github/workflows/cicd.yaml"

if [[ "${OLD_GOLANG_VERSION}" == "${NEW_GOLANG_VERSION}" ]]; then
    echo "No update needed, already on latest Golang version ${NEW_GOLANG_VERSION}"
    exit 0
fi

# we need to be at repo root to adjust go.mod
cd "${GIT_REPO_ROOT}" || exit 1

# shellcheck disable=SC2086
go mod edit -go=${NEW_GOLANG_VERSION}

# rename from $OLD_GOLANG_VERSION to $NEW_GOLANG_VERSION
# shellcheck disable=SC2116,SC2046,SC2086
grep -rFl "${OLD_GOLANG_VERSION}" $(echo "${FILES_NEEDING_UPDATES}") | ${XARGS_CMD} -r ${SED_CMD} "s/${OLD_GOLANG_VERSION//./[.]}/${NEW_GOLANG_VERSION}/g"
