#!/usr/bin/env bash

if [[ $# -ne 2 ]] ; then
  echo 'Usage: minio-bootstrap.sh <hostname> <bucket-name>'
  exit 1
fi

for i in 1 2 3 4 5 6; do
  echo 'Waiting for minio'
  curl -f http://$1:9000/minio/health/live
  ret=$?
  echo "Return code was $ret"
  [ $ret -eq 0 ] && break
  sleep 5
done

/usr/bin/mc alias set local http://$1:9000 minioadmin minioadmin;
/usr/bin/mc mb --ignore-existing local/$2;
/usr/bin/mc policy set public local/$2;
