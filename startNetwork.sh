#!/bin/bash

pushd network

./startnetwork.sh
sleep 5

./createchannel.sh
sleep5

./setAnchorPeerUpdate.sh

popd