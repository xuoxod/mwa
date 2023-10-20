#! /usr/bin/bash

clear

go build -o mwa ./cmd/web/*.go

./mwa
