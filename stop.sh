
#!/bin/sh

pid=`ps auxf | grep '\./floodproxy' | grep -v 'grep' | grep -v 'floodproxy' | awk '{print $2}'`

kill $pid
