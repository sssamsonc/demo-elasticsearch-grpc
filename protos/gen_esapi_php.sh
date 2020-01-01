#!/bin/sh

cd /var/www/protos

protoc --proto_path=./ --php_out=./ --grpc_out=./ --plugin=protoc-gen-grpc=/usr/bin/grpc_php_plugin ./esapi.proto

if [ -d /var/www/Esapi ]; then
  cp -r Esapi/* /var/www/Esapi/
  rm -rf Esapi
else
  mv -f Esapi /var/www/
fi

if [ -d /var/www/GPBMetadata ]; then
  cp -r GPBMetadata/* /var/www/GPBMetadata/
  rm -rf GPBMetadata
else
  mv -f GPBMetadata /var/www/
fi