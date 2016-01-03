# go-sphere

[![Build Status](https://travis-ci.org/samuelngs/go-sphere.svg?branch=master)](https://travis-ci.org/samuelngs/go-sphere)
[![Coverage Status](https://coveralls.io/repos/samuelngs/go-sphere/badge.svg?branch=master&service=github)](https://coveralls.io/github/samuelngs/go-sphere?branch=master)
[![GoDoc](https://godoc.org/github.com/samuelngs/go-sphere?status.svg)](https://godoc.org/github.com/samuelngs/go-sphere)

Go-Sphere is an open source realtime framework to scale websocket horizontally (across multiple hosts) via pub/sub  synchronization. Sphere supports both real-time bidirectional event-based communication and group communication.

## Usage

To create websocket server with gin webserver
```go
package main

import (
  "github.com/samuelngs/go-sphere"
  "github.com/gin-gonic/gin"
)

func main() {
    // create websocket server
    s := sphere.Default()
    // create http server
    r := gin.Default()
    // attach websocket handler
    r.GET("/sync", func(c *gin.Context) {
      s.Handler(c.Writer, c.Request)
    })
    // listen and serve on 0.0.0.0:8080
    r.Run(":8080") 
}
```

Scale websocket horizontally via pub/sub (broker) synchronization 
```go
b := sphere.NewRedisBroker()
s := sphere.Default(b) // <= pass in redis broker when creates websocket server
```

Use custom pubsub broker/agent
```go
type CustomBroker struct {
  *Broker
}

func (broker *CustomBroker) OnSubscribe(channel *Channel) error {
  return nil
}

func (broker *CustomBroker) OnUnsubscribe(channel *Channel) error {
  return nil
}

func (broker *CustomBroker) OnPublish(channel *Channel, data *Packet) error {
  return nil
}

func (broker *CustomBroker) OnMessage(channel *Channel, data *Packet) error {
  return nil
}

func main() {
  customBroker := &CustomBroker{ExtendBroker()}
  s := sphere.Default(customBroker)
}
```

Custom channel events
```go
type TestSphereModel struct {
  *ChannelModel
}

func (m *TestSphereModel) Subscribe(room string, connection *Connection) bool {
  // return false if you want to reject channel subscribe request
  return true
}

func (m *TestSphereModel) Disconnect(room string, connection *Connection) bool {
  return true
}

func (m *TestSphereModel) Receive(event string, message string) (string, error) {
  return "response_data", nil
}

func main() {
  s := sphere.Default()
  s.ChannelModels(&TestSphereModel{ExtendChannelModel("test")}) // <= set channel namespace "test"
}

```

## Documentation

`go doc` format documentation for this project can be viewed online without installing the package by using the GoDoc page at: https://godoc.org/github.com/samuelngs/go-sphere

## Contributing

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/samuelngs/go-sphere/issues)
- Fix bugs and [submit pull requests](https://github.com/samuelngs/go-sphere/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

## License ##

This project is distributed under the MIT license found in the [LICENSE](./LICENSE)
file.

```
The MIT License (MIT)

Copyright (c) 2015 Samuel

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
