#!/bin/sh

set -eo

go build -o ebi-reversi main.go && ./ebi-reversi
