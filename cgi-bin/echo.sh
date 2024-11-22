#!/bin/bash
echo "Content-Type: $CONTENT_TYPE"
printf "\r\n"
dd ibs=1 count=$CONTENT_LENGTH
exit 0