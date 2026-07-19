#!/bin/sh
set -eu

part="${1:-hotfix}"
version_file="${VERSION_FILE:-VERSION}"

if [ ! -f "$version_file" ]; then
  echo "missing $version_file" >&2
  exit 1
fi

version="$(tr -d '[:space:]' < "$version_file")"

case "$version" in
  [0-9]*.[0-9]*.[0-9]*) ;;
  *)
    echo "invalid version in $version_file: $version" >&2
    exit 1
    ;;
esac

old_ifs="$IFS"
IFS=.
set -- $version
IFS="$old_ifs"

if [ "$#" -ne 3 ]; then
  echo "invalid version in $version_file: $version" >&2
  exit 1
fi

major="$1"
minor="$2"
hotfix="$3"

case "$major$minor$hotfix" in
  *[!0-9]*)
    echo "invalid version in $version_file: $version" >&2
    exit 1
    ;;
esac

case "$part" in
  major)
    major=$((major + 1))
    minor=0
    hotfix=0
    ;;
  minor)
    minor=$((minor + 1))
    hotfix=0
    ;;
  hotfix|patch)
    hotfix=$((hotfix + 1))
    ;;
  *)
    echo "usage: $0 major|minor|hotfix" >&2
    exit 1
    ;;
esac

new_version="$major.$minor.$hotfix"
printf '%s\n' "$new_version" > "$version_file"
printf '%s\n' "$new_version"
