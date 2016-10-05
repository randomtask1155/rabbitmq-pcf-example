package main 

import ( 
  "github.com/streadway/amqp"
  "fmt"
  "crypto/tls"
)

func startConsumer(counter chan int64) {
  conn, err := amqp.DialTLS(amqpAddress, &tls.Config{InsecureSkipVerify: true})
  if err != nil {
    fmt.Printf("Could not connect to %s: %s\n", amqpAddress, err)
  }
  defer conn.Close()
  
  ch, err := conn.Channel()
  if err != nil {
    fmt.Printf("Could open channel: %s\n", err)
  }
  defer ch.Close()
  
  q, err := ch.QueueDeclare(
    "hello",
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )
  if err !=  nil {
    fmt.Printf("Could not declare Queue: %s\n", err)
  }
  
  msgs, err := ch.Consume(
                q.Name, // queue
                "",     // consumer
                true,   // auto-ack
                false,  // exclusive
                false,  // no-local
                false,  // no-wait
                nil,    // args
              )
  
  if err != nil {
    fmt.Printf("Consumer fialed: %s\n", err)
  }

  for range msgs {
    counter <- 1
  }
}