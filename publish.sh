#!/bin/sh
#last published version - v0.1.1
git tag $1 #vx.x.x
git push origin $1
GOPROXY=proxy.golang.org go list -m github.com/wizards-0/go-pins@$1