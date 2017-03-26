package main

import (
	"fmt"
	"time"
)

const (
	nwaiting   = 5
	nbarbers   = 2
	maxclients = 500
	sleep      = iota
	full
)

type Client struct {
	id int
	ch chan int
}

var (
	lobby           = make(chan *Client)
	waiting         = make(chan int)
	barbers         = make(chan int)
	quit            = make(chan int)
	clients_counter = 0
	clients_waiting = 0
	barbers_waiting = 0
)

func checkWaiting() bool {
	if clients_waiting < nwaiting {
		return true
	}

	return false
}

func attendClients() {
	if barbers_waiting > 0 && clients_waiting > 0 {
		barbers_waiting--
		clients_waiting--

		waiting <- 0
		barbers <- 0
	}
}

func receptionist() {
	for {
		select {
		case client := <-lobby:
			if checkWaiting() {
				clients_waiting++
				client.ch <- sleep
			} else {
				client.ch <- full
			}
		case <-quit:
			barbers_waiting++
		}

		attendClients()
	}
}

func launch_client(client *Client) {
	lobby <- client

	switch reply := <-client.ch; reply {
	case full:
		fmt.Printf("Cliente %d: me voy de la barbería, está llena\n", client.id)
		return
	case sleep:
		fmt.Printf("Cliente %d: me siento en la sala de espera\n", client.id)
		<-waiting
	}

	fmt.Printf("Cliente %d: me corto el pelo\n", client.id)
	fmt.Printf("Cliente %d: termino de cortarme el pelo\n", client.id)
}

func launch_barber(id int) {
	quit <- 0
	for {
		fmt.Printf("Barbero %d: me duermo esperando clientes\n", id)
		<-barbers
		fmt.Printf("Barbero %d: empiezo a cortar el pelo\n", id)
		fmt.Printf("Barbero %d: termino de cortar el pelo\n", id)
		quit <- 0
	}
}

func main() {
	go receptionist()

	for i := 0; i < nbarbers; i++ {
		go launch_barber(i + 1)
	}

	for clients_counter < maxclients {
		clients_counter++

		client := new(Client)
		client.id = clients_counter
		client.ch = make(chan int)

		go launch_client(client)
	}

	time.Sleep(5 * time.Second)
}
