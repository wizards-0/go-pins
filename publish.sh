#!/bin/sh
git tag $1 #vx.x.x
git push origin $1
GOPROXY=proxy.golang.org go list -m github.com/wizards-0/go-pins@$1