#!/bin/sh
set -eu

jq -n \
  --arg version "${1:?usage: release-metadata.sh VERSION COMMIT DATE}" \
  --arg commit "${2:?missing commit}" \
  --arg date "${3:?missing date}" \
  '{version: $version, commit: $commit, date: $date}' > version.json
