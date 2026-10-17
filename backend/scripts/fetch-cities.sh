#!/usr/bin/env sh
set -eu

script_directory=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
catalog_directory="$script_directory/../data"
mkdir -p "$catalog_directory"

curl --fail --show-error --location --retry 2 --max-time 30 \
  https://raw.githubusercontent.com/hflabs/city/master/city.csv \
  --output "$catalog_directory/city.csv.download"
mv "$catalog_directory/city.csv.download" "$catalog_directory/city.csv"
