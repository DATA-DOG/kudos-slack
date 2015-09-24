#!/bin/bash
#darwin
#386
docker run --rm -it -v "$PWD":/usr/src/myapp -v "$GOPATH":/go -w /usr/src/myapp golang:1.4.2-cross bash -c '
for GOOS in  linux; do
  for GOARCH in amd64; do
    go build -v -o kudos-$GOOS-$GOARCH
  done
done
'
