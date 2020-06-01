#!/bin/bash

basepath=$(cd `dirname $0`; pwd)

protoc -I $basepath --go_out=$basepath $basepath/cs.proto
