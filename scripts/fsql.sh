#!/bin/bash
creekGenUrl='http://creek.baidubce.com/v1/creek/generate?flag=exe&architecture=linux-amd64'
if [ "$1" == "" ] || [ "$2" == ""  ] || [ "$3" == "" ];then
    echo "usage: fsql.sh <job.json> <file> <sql>"
    exit 0
fi
 
if [[ ! -f "$1" ]];then
  echo "oops! file does not found: $1"
fi
 
json=$(cat "$1")
json="${json/FILE_PLACE_HOLDER/$2}"
json="${json/SQL_PLACE_HOLDER/$3}"

#exit 0
code=$(curl --compressed -l -k -H 'Content-Type:application/json;charset=utf-8' -o creek -w %{http_code} -X POST -d "$json" $creekGenUrl)
if [[ $code -eq "200" ]];then
         chmod 775 creek
else
    echo "failed to generate creek, http code:$code"    
fi

./creek
