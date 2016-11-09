#!/bin/bash

loglevel=$1
file=$2

sed -ie "s/{LOGLEVEL}/$loglevel/g" $file