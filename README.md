# rabbitmq-pcf-example

## producer.go

Creates a rabbit queue called "hello" and publishes "hello" ever 1 second to the queue

## consumer.go

Connects to the "hello" queue and consumes the "hello" messages

## main.go

Sets up the settings for the producer and consumer before launching them.  Also starts the http server that serves the web page and /getcounts endpoint

Keeps track of all messages published by producer and consumed by consumer.  The counters will rotate every 100,000 messages and will be displayed via the web interface

# Build and run in PCF
```
git clone git@github.com:randomtask1155/rabbitmq-pcf-example.git
cd rabbitmq-pcf-example
cf push
```