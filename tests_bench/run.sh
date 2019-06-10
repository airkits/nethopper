#!/bin/bash

basepath=$(cd `dirname $0`; pwd)

go test -bench=. $basepath/log/*