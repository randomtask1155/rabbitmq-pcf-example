package main 

import (
  "encoding/json"
  "fmt"
  "os"
  "net/http"
  "html/template"
)

var (
  amqpAddress = "amqp://guest:guest@localhost:5672/"
  producerCount int64 
  consumerCount int64
  
  startPageTemplate = template.Must(template.ParseFiles("tmpl/start.tmpl")) // root page
)


/*
  Fetach AMQP service address from VCAP services
  "VCAP_SERVICES": {
      "p-rabbitmq": [
        {
          "credentials": {
            "http_api_uris": [
              "https://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@pivotal-rabbitmq.run-07.haas-59.pez.pivotal.io/api/"
            ],
            "ssl": true,
            "dashboard_url": "https://pivotal-rabbitmq.run-07.haas-59.pez.pivotal.io/#/login/eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65/ihgpjelrhbocd6pd0sj3lo8ii5",
            "password": "ihgpjelrhbocd6pd0sj3lo8ii5",
            "protocols": {
              "management+ssl": {
                "path": "/api/",
                "ssl": false,
                "hosts": [
                  "10.193.72.76"
                ],
                "password": "ihgpjelrhbocd6pd0sj3lo8ii5",
                "username": "eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65",
                "port": 15672,
                "host": "10.193.72.76",
                "uri": "http://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@10.193.72.76:15672/api/",
                "uris": [
                  "http://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@10.193.72.76:15672/api/"
                ]
              },
              "amqp+ssl": {
                "vhost": "2c03cdfd-1c7d-4bf0-9a4d-fcf7861071d9",
                "username": "eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65",
                "password": "ihgpjelrhbocd6pd0sj3lo8ii5",
                "port": 5671,
                "host": "10.193.72.76",
                "hosts": [
                  "10.193.72.76"
                ],
                "ssl": true,
                "uri": "amqps://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@10.193.72.76:5671/2c03cdfd-1c7d-4bf0-9a4d-fcf7861071d9",
                "uris": [
                  "amqps://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@10.193.72.76:5671/2c03cdfd-1c7d-4bf0-9a4d-fcf7861071d9"
                ]
              }
            },
            "username": "eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65",
            "hostname": "10.193.72.76",
            "hostnames": [
              "10.193.72.76"
            ],
            "vhost": "2c03cdfd-1c7d-4bf0-9a4d-fcf7861071d9",
            "http_api_uri": "https://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@pivotal-rabbitmq.run-07.haas-59.pez.pivotal.io/api/",
            "uri": "amqps://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@10.193.72.76/2c03cdfd-1c7d-4bf0-9a4d-fcf7861071d9",
            "uris": [
              "amqps://eb60ebbd-fb8e-4e02-8f2a-a73a8f664b65:ihgpjelrhbocd6pd0sj3lo8ii5@10.193.72.76/2c03cdfd-1c7d-4bf0-9a4d-fcf7861071d9"
            ]
          },
          "syslog_drain_url": null,
          "label": "p-rabbitmq",
          "provider": null,
          "plan": "standard",
          "name": "danl-rabbit",
          "tags": [
            "rabbitmq",
            "messaging",
            "message-queue",
            "amqp",
            "stomp",
            "mqtt",
            "pivotal"
          ]
        }
      ]
    }
*/
func setAMQPServiceAddress() {
    vcap := os.Getenv("VCAP_SERVICES")
    if vcap == "" {
      fmt.Printf("Using default AMQP address '%s'\n", amqpAddress)
      return
    }
    
    type Creds struct {
      URI string `json:"uri"`
    }
    type PRabbit struct {
        Credentials Creds `json:"credentials"`
    } 
    type VcapServices struct {
      Rabbit []PRabbit `json:"p-rabbitmq"`
    }
    
    vs := new(VcapServices)
    err := json.Unmarshal([]byte(vcap),&vs)
    if err != nil {
      fmt.Printf("Unable to parse VCAP_SERVICES: %s\nERORR:%s\n", vcap, err)
    }
    
    for i := range vs.Rabbit {
      if vs.Rabbit[i].Credentials.URI != "" {
        amqpAddress = vs.Rabbit[i].Credentials.URI
        fmt.Printf("Using AMQP address '%s'\n", amqpAddress)
        return
      }
    }
    // if we made it here then just print the default string
    fmt.Printf("Using default AMQP address '%s'\n", amqpAddress)
}

func runHTTPServer(port string) {
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/getcounts", getCounts)
  http.Handle("/img/", http.FileServer(http.Dir("")))
	http.Handle("/js/", http.FileServer(http.Dir("")))
	http.Handle("/css/", http.FileServer(http.Dir("")))  
  err := http.ListenAndServe(":"+port, nil)
  if err != nil {
    fmt.Printf("Failed to start HTTP Server: %s\n", err)
    os.Exit(2)
  }
}

/*
  server the root / page 
*/
func rootHandler(w http.ResponseWriter, r *http.Request) {
  startPageTemplate.Execute(w, "")
}

/*
  Browser will call /getcounts and this function will return 
  producter and consumer counts
*/
func getCounts(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  
  type Resp struct  {
    ProducerCount int64 `json:"producercount"`
    ConsumerCount int64 `json:"consumercount"`
  }
  resp := Resp{producerCount,consumerCount}
  jdata, err := json.Marshal(resp)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    emsg, _ := buildJSONError(fmt.Sprintf("Failed to marshal data: %s\n", err))
    w.Write(emsg)
  }
  w.Write(jdata)
}

/*
  build a json formatted error for http response
*/
func buildJSONError(msg string) ([]byte, error) {
  type JError struct {
    Error string `json:"error"`
  }
  je := JError{msg}
  return json.Marshal(je)
}

/*
   Use functions to update the producer and consumer counts so we can create mutex lock 
   given these will be updated by go routines.   Not required for this scope but 
   implemented for example purposes.
*/
func updateCount(counter *int64, c int64) {
  if *counter >= 100000 {
    *counter = 0
    return
  }
  *counter += c
}

/*
  Setup the http server and start the RabbitMQ producers and consumers
*/
func main() {
  setAMQPServiceAddress()
  port := os.Getenv("PORT")
  if port == "" {
    fmt.Printf("Environment variable $PORT is not set\n")
    os.Exit(1)
  }
  
  producerChan := make(chan int64,0)
  consumerChan := make(chan int64,0)
  
  go runHTTPServer(port)
  go startProducer(producerChan)
  go startConsumer(consumerChan)
  
  for {
      select {
      case p := <-producerChan:
        updateCount(&producerCount, p)
      case c := <-consumerChan:
        updateCount(&consumerCount, c)
      }
  }

}