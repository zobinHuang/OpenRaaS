#!/usr/bin/env bash

cd ../dockertools

sh ./rm-wine.sh 1

umount ../winetools/apps/point*
rm -rf ../winetools/apps/point*