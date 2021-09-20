#!/bin/sh

set -eu

go build -o ebi-reversi main.go && ./ebi-reversi
