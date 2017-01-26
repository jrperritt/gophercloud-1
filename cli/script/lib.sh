#!/bin/bash

# This script is intended to be `source`d

# It populates several environment variables via these functions:
# $BRANCH - get_branch()
# $VERSION - get_version()
# $COMMIT - get_commit()

get_branch() {
  # See http://docs.travis-ci.com/user/environment-variables/#Default-Environment-Variables
  # for details about default Travis Environment Variables and their values
  if [ -z "${TRAVIS_BRANCH-}" ]; then
    BRANCH=$(git rev-parse --abbrev-ref HEAD)
  else
    BRANCH=${TRAVIS_BRANCH}
  fi
  export BRANCH

  return 0
}

get_commit() {
  COMMIT=$(git rev-parse --verify HEAD)
  RETURN_CODE=$?
  export COMMIT
  return $RETURN_CODE
}

get_version() {
  # Version will be the most recent tag + "-dev" if working tree is dirty
  # e.g.: "1.0.1"" or "1.0.1-dev"
  VERSION=$(git describe --tags --dirty='-dev' 2> /dev/null)
  export VERSION
  return 0
}

#
# Helper function to do replace; this should work across operating systems
#
update() {
  TMP_FILE=$(mktemp "$1")
  sed -e "$2" "$3" > "$TMP_FILE"
  chmod 0644 "$TMP_FILE"
  mv -f "$TMP_FILE" "$3"
}
