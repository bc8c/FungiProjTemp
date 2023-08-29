#!/bin/bash

pushd application

    ./public/ccp/ccp-generate.sh
    sleep 2

    npm install
    sleep 10

    npm start

popd