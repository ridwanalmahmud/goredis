package main

import (
	"fmt"
	"log"
	"sync"
	"context"
	"time"
	"testing"
	"github.com/Ridwan-Al-Mahmud/goredis/client"
)

func TestServerWithMultiClients(t *testing.T) {
	server := NewServer(Config{})
	go func(){
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)
	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func(it int){
			c, err := client.New("localhost:5001")
			if err != nil {
				log.Fatal(err)
			}
			defer c.Close()
			key := fmt.Sprintf("client_foo_%d", i)
			value := fmt.Sprintf("client_bar_%d", i)
			if err := c.Set(context.TODO(), key, value); err != nil {
				log.Fatal(err)
			}
			val, err := c.Get(context.TODO(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Client %d got this val back => %s\n", it, val)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)
	if len(server.peers) != 0 {
		t.Fatalf("Server should have no peers, but has %d", len(server.peers))
	}
}