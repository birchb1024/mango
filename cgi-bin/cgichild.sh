#!/bin/bash
echo "Content-type: text/plain"
#echo "Content-length: $CONTENT_LENGTH"
echo
echo
echo $CONTENT_LENGTH
env
dd ibs=1 count=$CONTENT_LENGTH
exit 0