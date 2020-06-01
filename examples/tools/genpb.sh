#!/bin/bash

basepath=$(cd `dirname $0`; pwd)
$basepath/../../network/transport/pb/cs/gen.sh
$basepath/../../network/transport/pb/ss/gen.sh
$basepath/../model/pb/c2s/gen.sh
$basepath/../model/pb/s2s/gen.sh
