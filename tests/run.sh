#!/bin/bash

basepath=$(cd `dirname $0`; pwd)
go test -v $basepath/log/*.go
go test -v $basepath/server/*.go
go test -v $basepath/utils/*.go
go test -v $basepath/cache/*.go
go test -v $basepath/database/*.go
go test -v $basepath/network/*.go
go test -v $basepath/timer/*.go