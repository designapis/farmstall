#!/usr/bin/env bash


FQDN=https://farmstall.example.com

docker stop farmstall || echo

docker run -d --name farmstall -it --rm -p 9999:9999 -e FQDN=$FQDN -e PORT=9999 designapis/farmstall:latest

FQDN=$FQDN URL=http://localhost:9999 strest
testCode=$?

docker stop farmstall

exit $testCode
