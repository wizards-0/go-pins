#!/bin/sh

mainVersionLine=$(grep "main.version" version.ini)
mainVersionKeyValue=(${mainVersionLine//=/ })
mainVersion=${mainVersionKeyValue[1]}
mainVersionParts=(${mainVersion//./ })
mainVersionPatchOld=${mainVersionParts[2]}
mainVersionPatchNew=$((mainVersionPatchOld+1))
mainVersionNew="${mainVersionParts[0]}.${mainVersionParts[1]}.$mainVersionPatchNew"
mainVersionLineNew="main.version=$mainVersionNew"
sed -i "s/${mainVersionLine}/${mainVersionLineNew}/g" ./version.ini
git commit -a -u -m "Automated commit for version increment"
git tag $mainVersionNew
git push origin $mainVersionNew
GOPROXY=proxy.golang.org go list -m github.com/wizards-0/go-pins@$mainVersionNew
