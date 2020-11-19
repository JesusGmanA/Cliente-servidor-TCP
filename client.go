package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/exec"

	"./process"
	"./useful"
)

const ADD_CLIENT = 1
const DELETE_CLIENT = 2
const EXIT = 3

type Client struct {
	ClientId int64
	Process  process.Process
}

func (c *Client) createClient() {
	var cli Client
	client, err := net.Dial("tcp", useful.PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(client).Encode(c) //Initial request value the server responds with a clientId as well as a Process
	err = gob.NewDecoder(client).Decode(&cli)
	if err != nil {
		fmt.Println(err)
		return
	}
	*c = cli
	fmt.Print(c.Process)
	c.Process.ContinueRunning = true
	c.Process.StartProcess()
	client.Close()
}

func (c *Client) deleteClient() {
	client, err := net.Dial("tcp", useful.PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Process.ContinueRunning = false
	err = gob.NewEncoder(client).Encode(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Client has stopped")
	client.Close()
}

func main() {
	c := Client{ClientId: -1, Process: process.Process{}}
	go c.createClient()
	fmt.Print("Press 'Enter' to exit...")
	fmt.Scanln() //El primero se come el "Enter" atorado en el buffer cuando se lee algo desde consola
	c.deleteClient()
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
