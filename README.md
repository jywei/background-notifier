`redis-cli -r 100 -p 4000 RPUSH resque:queue:slack '{"class":"notifier","args":["This is a test notification"]}'`

`echo "port 4000" | redis-server -`
