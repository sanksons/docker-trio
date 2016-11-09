#!/bin/bash

loglevel=$1
file=$2

sed -i "s/{LOGLEVEL}/$loglevel/g" $file