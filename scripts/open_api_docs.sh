#!/bin/bash
## http://localhost:5042/swagger/v1/index.html#/

url="http://localhost:5042/swagger/v1/index.html#/"

name=$(uname -s)

if [[ "$name" == "Linux" ]]; then
    xdg-open $url
elif [[ "$name" == "Darwin" ]]; then
    open $url
fi
