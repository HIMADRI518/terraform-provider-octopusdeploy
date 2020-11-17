#!/bin/bash

# Dot Source this script to set the RELEASE_VERSION environment variable

set -eo pipefail

# if [ -z "${CIRCLE_BUILD_NUM}" ]; then
#     echo "The environment variable CIRCLE_BUILD_NUM is not set. Setting as 999."
#     CIRCLE_BUILD_NUM="999"
# fi

# RELEASE_VERSION="0.0.2-alpha.${CIRCLE_BUILD_NUM}"
# echo "Release version is ${RELEASE_VERSION}"

# export RELEASE_VERSION=$RELEASE_VERSION

# Get latest tag ref
REF=$(git rev-list --tags --max-count=1 HEAD)

# Get latest tag 
LATEST_TAG=$(git describe --tags ${REF}) 

# Remove suffix
LATEST_TAG=${LATEST_TAG/v/''}

# Determine next version
RELEASE_VERSION=`echo $LATEST_TAG | awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}'`

# Add suffix
RELEASE_VERSION="v${RELEASE_VERSION}"

echo "Release version is ${RELEASE_VERSION}"
export RELEASE_VERSION=$RELEASE_VERSION