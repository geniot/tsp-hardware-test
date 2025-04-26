#!/bin/bash

pkill -f hwt

if [ "$#" -gt 0 ]; then
 /mnt/SDCARD/Apps/HardwareTest/hwt "$@"
else
  progdir=$(dirname "$0")
  cd $progdir
  ./hwt
fi