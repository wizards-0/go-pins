#!/bin/sh
#last published version
#git tag $1 #vx.x.x
#git push origin $1
#GOPROXY=proxy.golang.org go list -m github.com/wizards-0/go-pins@$1
mainVersionLine=$(grep "main.version" version.ini)
mainVersionKeyValue=(${mainVersionLine//=/ })
mainVersion=${mainVersionKeyValue[1]}
mainVersionParts=(${mainVersion//./ })
mainVersionPatchOld=${mainVersionParts[2]}
mainVersionPatchNew=$((mainVersionPatchOld+1))
mainVersionLineNew="main.version=${mainVersionParts[0]}.${mainVersionParts[1]}.$mainVersionPatchNew"
sed -i "s/${mainVersionLine}/${mainVersionLineNew}/g" ./version.ini
cat version.ini