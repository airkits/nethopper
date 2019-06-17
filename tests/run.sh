#!/bin/bash

basepath=$(cd `dirname $0`; pwd)
go test $basepath/log/*.go
go test $basepath/server/*.go
go test $basepath/utils/*.go