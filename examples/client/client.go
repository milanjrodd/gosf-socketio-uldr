package main

import (
	"log"
	"runtime"
	"time"

	gosocketio "github.com/milanjrodd/gosf-socketio-uldr"
	"github.com/milanjrodd/gosf-socketio-uldr/transport"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < 100; i++ {

		go func() {
			c, err := gosocketio.Dial(
				gosocketio.GetUrl("localhost", 3811, false),
				transport.GetDefaultWebsocketTransport(),
				nil,
			)
			if err != nil {
				log.Fatal(err)
			}

			err = c.On("/message", func(h *gosocketio.Channel, args Message) {
				log.Println("--- Got chat message: ", args)
			})
			if err != nil {
				log.Fatal(err)
			}

			err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
				log.Fatal("Disconnected")
			})
			if err != nil {
				log.Fatal(err)
			}

			err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
				log.Println("Connected")
			})
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(1 * time.Second)

			for i := 0; i < 100; i++ {
				go sendJoin(c)
				go sendJoin(c)
				go sendJoin(c)
				go sendJoin(c)
				go sendJoin(c)
			}

			time.Sleep(10 * time.Second)
			c.Close()

		}()
	}
	time.Sleep(6000 * time.Second)
	log.Println(" [x] Complete")
}
