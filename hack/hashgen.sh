#!/bin/sh

for f in bin/jaas*; do shasum -a 256 $f > $f.sha256; done

