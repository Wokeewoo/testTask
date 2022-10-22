package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"time"
)

func connectToNats(natsUrl string) (*nats.Conn, error) {
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return nil, err
	}
	return nc, nil
}

func connectToStan(clusterID string, clientID string) (stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		return nil, err
	}
	return sc, nil
}
func GetChanMsgs() {
	nc, err := connectToNats(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	sc, err := connectToStan("test-cluster", "Aba")
	if err != nil {
		log.Fatal(err)
	}
	sub, err := subscribeToChannel(sc)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to
			sub.Close()
			sc.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}

func subscribeToChannel(sc stan.Conn) (stan.Subscription, error) {
	sub, err := sc.Subscribe("foo", myFunc, stan.DeliverAllAvailable())
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func myFunc(msg *stan.Msg) {
	fmt.Printf("sss %s\n", msg.Data)
	var model Model
	err := json.Unmarshal(msg.Data, &model)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse msg to a data model: %v\n", err)
		log.Fatal(err)
	}
	cache.Set(model.OrderUid, &model, 5*time.Hour)

	// Получить кеш с ключем "myKey"
	i, _ := cache.Get(model.OrderUid)
	b, _ := json.Marshal(i)
	fmt.Printf("vrot mne nogiiiiiii    \n %s", b)
}
