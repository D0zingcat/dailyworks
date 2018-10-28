#!/bin/bash
docker run -d --rm --name ss -v /home/d0zingcat/dailyworks/ss/config.json:/root/config.json  --name ss -p 19841:19841/tcp -p 11235:11235/tcp -p 12580:12580/tcp -p 13311:13311/tcp -p 15151:15151/tcp -p 53211:53211/tcp \
       	-p 19999:19999/tcp -p 22222:22222/tcp \
-p 12121:12121/tcp -p 34555:34555/tcp -p 21234:21234/tcp -p 19876:19876/tcp -p 11724-11736:11724-11736/tcp -p 11724-11736:11724-11736/udp d0zingcat/ss-goplus -u -c /root/config.json 
