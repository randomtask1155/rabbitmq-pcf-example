package main 

import (
  "github.com/streadway/amqp"
  "fmt"
  "time"
  "crypto/tls"
)

func startProducer(counter chan int64) { 
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
    fmt.Printf("Could not create Queue: %s\n", err)
  }
  
  body := "hello"
  for {
    err = ch.Publish(
      "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
          ContentType: "text/plain",
          Body:        []byte(body),
        })
    if err != nil {
      fmt.Printf("Could not publish record: %s\n", err)
    }
    counter <- 1 // update counters 
    time.Sleep(1 * time.Second) // send message every 1 second
  }
  
}