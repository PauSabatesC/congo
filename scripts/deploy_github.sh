#!/bin/bash

version=$1 #e.g. v0.1.1

if [[ -z $version ]]; then
  echo "ERROR: specify a version e.g v0.1.1"
  exit 1
fi

read -p "Did you push the local commit? (y/n)" -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 0
fi

echo "Creating TAG $version ..."
sleep 5
git tag -a $version -m "Version $version"
git push origin $version

echo "A github action is triggered to build the release for $version! Great job!"
