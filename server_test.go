package main

import (
	"fmt"
	"log"
	"context"
	"time"
	"testing"
	"github.com/redis/go-redis/v9"
)

func TestRedisClient(t *testing.T) {	
	listenAddr := ":5001"
	server := NewServer(Config{
		ListenAddr: listenAddr,
	})
	go func(){
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Millisecond * 400)
	rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("localhost%s", ":5001"),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

	testCases := map[string]string {
		"foo": "bar",
		"game": "football",
		"score": "4-1",
		"team": "barca",
	}

	for key, val := range testCases {
	if err := rdb.Set(context.Background(), key, val, 0).Err(); err != nil {
        t.Fatal(err)
    }

    newVal, err := rdb.Get(context.Background(), key).Result()
    if err != nil {
        t.Fatal(err)
	}
	if newVal != val {
		t.Fatalf("expected %s got %s", val, newVal)
	}
	fmt.Printf("key: %s => val: %s\n", key, newVal)

	} 
}
/*
func TestFooBar(t *testing.T) {
	in := map[string]string{
		"server": "redis",
	}
	out := respWriteMap(in)
	fmt.Println(out)
}

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
}*/ 
