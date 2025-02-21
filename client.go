package gosocketio

import (
	"net"
	"strconv"

	"github.com/milanjrodd/gosf-socketio-uldr/transport"
)

const (
	webSocketProtocol       = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl             = "/socket.io/?EIO=3&transport=websocket"
)

/*
*
Socket.io client representation
*/
type Client struct {
	methods
	Channel
	transport.Transport
	url string
}

/*
*
Get ws/wss url by host and port
*/
func GetUrl(host string, port int, secure bool) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	return prefix + net.JoinHostPort(host, strconv.Itoa(port)) + socketioUrl
}

/*
*
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(url string, tr transport.Transport, c *Client) (*Client, error) {
	if c == nil {
		c = &Client{}
		c.initMethods()
	}

	c.initChannel()

	c.Transport = tr
	c.url = url

	var err error
	c.conn, err = tr.Connect(url)
	if err != nil {
		return nil, err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)
	go heartbeat(&c.Channel)

	return c, nil
}

/*
*
Close client connection
*/
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}

/*
*
Open client connection
*/
func (c *Client) Open() (err error) {
	c.initChannel()

	c.conn, err = c.Transport.Connect(c.url)

	if err != nil {
		return err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)
	go heartbeat(&c.Channel)

	return nil
}
