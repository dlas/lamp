
#!/bin/bash

cd `dirname $0`

export WEBROOT=`pwd`/webdata

bin/webmain

i2cset -y 2 65 51 0

