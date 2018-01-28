# Open another port for this background worker
`echo "port 4000" | redis-server -`

# Send builk messages into redis
`redis-cli -r 100 -p 4000 RPUSH resque:queue:slack '{"class":"notifier","args":["This is a test notification"]}'`

# Run the service
`go run main.go -queues=notifier`