#!/bin/bash

WD=$(pwd)
WORKDIR=tmp
TFVERSION=0.7.11
if [[ $OSTYPE =~ darwin ]]; then
  TFARCH=darwin_amd64
else
  TFARCH=linux_amd64
fi
TFURL="https://releases.hashicorp.com/terraform/${TFVERSION}/terraform_${TFVERSION}_${TFARCH}.zip"

if ! [[ -d $WORKDIR ]]; then
  mkdir $WORKDIR
fi
cd $WORKDIR

if ! [[ -e terraform ]]; then
  echo "downloading terraform"
  curl -s $TFURL -o terraform_${TFVERSION}_${TFARCH}.zip
  if [[ $? -ne 0 ]]; then
    echo "failed to download terraform"
    exit 1
  fi
  unzip terraform_${TFVERSION}_${TFARCH}.zip
fi

cd "$WD"

if [[ "$DEBUG" != "true" ]]; then
  echo "Downloading GO dependencies"
  go get -v
  if [[ $? -ne 0 ]]; then
    echo "Failed to download all dependencies"
    exit 1
  fi
fi

echo "Building terraform-provider-etcd:"
go build -v
if [[ $? -ne 0 ]]; then
  echo "Failed to build terraform-provider-etcd"
  exit 1
fi
cp terraform-provider-etcd $WORKDIR
cp testing/* $WORKDIR
cd $WORKDIR

if ! grep "${WD}/${WORKDIR}/test/terraform-provider-etcd" ~/.terraformrc 2>&1 > /dev/null; then
  echo
  echo "You'll have to change your ~/.terraformrc file to include this"
  echo "if you want to continue running these tests:"
  echo
  echo "providers {"
  echo "  etcd = \"${WD}/${WORKDIR}/test/terraform-provider-etcd\""
  echo "}"
  exit 1
fi

# cleanup just in case
docker stop $(docker-compose ps -q) &>/dev/null
docker-compose kill &>/dev/null
docker-compose rm &>/dev/null

echo "Setting up Etcd"
docker-compose run -p 2379:2379 -d etcd >/dev/null

if [[ $OSTYPE =~ darwin ]]; then
  if [ -z "$DOCKER_HOST" ]; then
    echo "Missing DOCKER_HOST environment key"
    exit 1
  fi
  hostport=${DOCKER_HOST##*//}
  ETCD_AUTHORITY="${hostport%%:*}:2379"
else
  ETCD_AUTHORITY="$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker-compose ps -q)):2379"
fi

if [[ "$ETCD_AUTHORITY" == "" ]]; then
  echo "Failed to get Etcd endpoint"
  exit 1
fi
echo "(ETCD_AUTHORITY=$ETCD_AUTHORITY)"

sleep 5s

rm -rf test; mkdir test
sed "s/PLACEHOLDER/${ETCD_AUTHORITY}/" provider.tf > test/provider.tf

cp terraform test/
cp terraform-provider-etcd test/

echo
echo "Testing:"
RESOURCES="${TESTS:-simple}"
for i in $RESOURCES; do
  tffile="${WD}/testing/test_${i}.tf"
  if [[ -e $tffile ]]; then
    cp "$tffile" test/
    cd test

    if [[ "$DEBUG" == "true" ]]; then
      TF_LOG=DEBUG ./terraform apply
    else
      RES="$(./terraform apply >/dev/null)"
    fi
    if [[ $? -ne 0 ]]; then
      echo "$RES"
      echo "Failed to terraform apply (${tffile})"
      exit 1
    else
      echo "${i}: OK"
    fi
    rm "test_${i}.tf"
  else
    echo "${i} - Not implemented"
  fi
done
