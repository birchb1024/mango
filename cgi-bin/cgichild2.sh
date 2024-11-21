#!/bin/bash
echo "Content-type: text/html"
echo ""
echo '<html><body>'
echo 'Hello From Bash <br/>Environment:'
echo '<pre>'
/usr/bin/env
echo '</pre>'
echo '</body></html>'
exit 0