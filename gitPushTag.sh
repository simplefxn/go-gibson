#!/bin/bash
MESSAGE=$1
VERSION=$2
git add --all .
git commit -m "${MESSAGE}"
git tag -a ${VERSION} -m "${MESSAGE}"
git push origin --tags