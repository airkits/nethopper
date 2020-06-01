#!/bin/bash

basepath=$(cd `dirname $0`; pwd)

protoc -I $basepath --go_out=plugins=grpc:$basepath $basepath/s2s.proto