#!/bin/bash

set -e

# Find the directory we exist within
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}/..

BUILDDIR=$(pwd)/build

# Make dir
mkdir -p $BUILDDIR

# Clean build bin dir
rm -rf $BUILDDIR/*

function fail () {
	echo "Aborting due to failure." >&2
	exit 2
}

# Build binary
cd src/cmd
for bin in *; do
  cd $bin
  set -x
  go build -o $BUILDDIR/$bin || fail
  set +x
  cd -
done
