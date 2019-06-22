#!/bin/bash

basepath=$(cd `dirname $0`; pwd)
go test -v $basepath/log/*.go
go test -v $basepath/server/*.go
go test -v $basepath/utils/*.go