#!/bin/sh

set -eu

go build -o ebi-reversi && ./ebi-reversi
