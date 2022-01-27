#!/bin/bash

basepath=$(cd `dirname $0`; pwd)
$basepath/../../network/transport/pb/cs/gen.sh
$basepath/../../proto/ss/gen.sh
$basepath/../model/pb/c2s/gen.sh
$basepath/../model/pb/s2s/gen.sh
