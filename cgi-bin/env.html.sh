#!/bin/bash
echo "Content-type: text/html"
echo ""
echo '<html><body>'
echo ' <h1>CGI Environment</h1>'
echo '<pre>'
/usr/bin/env | sort
echo '</pre>'
echo '</body></html>'
exit 0