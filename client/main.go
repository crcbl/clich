package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	input = make(chan string)
	done  = make(chan struct{})
)

func main() {
    // Some shells aren't enabled to read newlines properly
	_ = exec.Command("stty", "sane")

    // Signal interrupts to close the connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

    c, err := connect()
    if err != nil {
        log.Fatal(err)
        return
    }
	defer c.Close()

	go write(c)
	go read()

	for {
		select {
		case msg := <-input:
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

            // TODO: establish a handshake to close the connection
			err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func connect() (*websocket.Conn, error) {
	// TODO: get the server's configuration
	port := "8080"
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:" + port, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dialer error:", err)
        return nil, err
	}

    return c, nil
}

func write(c *websocket.Conn) {
	defer close(done)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("%s", msg)
	}
}

func read() {
    // TODO: create bounds for the user name for format's sake
	println("Enter session user name:")
	reader := bufio.NewReader(os.Stdin)
	user, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

    // TODO: create bounds for message size based on buffer size
	fmt.Println("Enter text and press enter to send a message")
	defer close(input)
	for {
		msg, err := reader.ReadString('\n')
        clearIn()
		if err != nil {
			return
		}
		// TODO: allow user to configure a tidy mode (strip newlines)
		input <- fmt.Sprintf("%s-- %s", user, string(msg))
	}
}

// Clears the previously entered text
// TODO: can we get rid of the new line this leaves behind?
func clearIn() {
	println("\033[1A" + "\033[K")
}
