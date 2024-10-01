#!/usr/bin/env bash
set -Eeuo pipefail

if ! command -v go; then
    echo "Golang not installed, exiting"
    exit 1
fi

OLD_GOLANG_VERSION="1.23"
NEW_GOLANG_VERSION="${1}"
GIT_REPO_ROOT=$(git rev-parse --show-toplevel)
FILES_NEEDING_UPDATES="${GIT_REPO_ROOT}/Dockerfile* ${GIT_REPO_ROOT}/README.md ${GIT_REPO_ROOT}/scripts/update_golang_version.sh ${GIT_REPO_ROOT}/.github/workflows/*"

# we need to be at repo root to adjust go.mod
cd "${GIT_REPO_ROOT}" || exit 1

# shellcheck disable=SC2086
go mod edit -go=${NEW_GOLANG_VERSION}
#go mod tidy

# rename from $OLD_GOLANG_VERSION to $NEW_GOLANG_VERSION
# shellcheck disable=SC2116,SC2046
grep -rl "${OLD_GOLANG_VERSION}" $(echo "${FILES_NEEDING_UPDATES}") | xargs sed -i "" -e "s/${OLD_GOLANG_VERSION}/${NEW_GOLANG_VERSION}/g"

# show diff output so user can verify their changes
git diff
