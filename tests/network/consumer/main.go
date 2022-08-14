package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	StreamName          = "vehicle_data"
	LocationSubject     = "location"
	BatteryLevelSubject = "battery_level"
)

type (
	EventFields struct {
		Longitude    float64 `json:"longitude,omitempty"`
		Latitude     float64 `json:"latitude,omitempty"`
		BatteryLevel int     `json:"batteryLevel,omitempty"`
	}

	Event struct {
		MessageType string      `json:"msgType"`
		Timestamp   int64       `json:"timestamp"`
		Fields      EventFields `json:"fields"`
	}
)

func CreateStreamIfNotExists(streamName string, jc nats.JetStreamContext) error {
	_, err := jc.StreamInfo(streamName)
	if err != nil && !errors.Is(err, nats.ErrStreamNotFound) {
		return fmt.Errorf("failed to check if stream exist: %w", err)
	}
	jc.DeleteStream(streamName)
	if errors.Is(err, nats.ErrStreamNotFound) {
		_, err = jc.AddStream(&nats.StreamConfig{
			Name:   streamName,
			MaxAge: 0,
			Subjects: []string{
				streamName + ".*",
			},
		})
		if err != nil {
			return fmt.Errorf("failed to create stream: %w", err)
		}
	}

	return err
}

func PublishMessages(jc nats.JetStreamContext) error {
	var (
		ev    Event
		bytes []byte
		err   error
	)

	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < 10; i++ {
		ev = Event{
			MessageType: LocationSubject,
			Timestamp:   time.Now().UnixNano(),
			Fields: EventFields{
				Longitude: float64(rand.Int63n(360) - 180),
				Latitude:  float64(rand.Int63n(180) - 90),
			},
		}
		if bytes, err = json.Marshal(&ev); err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}
		if _, err = jc.PublishAsync(StreamName+"."+LocationSubject, bytes); err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}

		ev = Event{
			MessageType: BatteryLevelSubject,
			Timestamp:   time.Now().UnixNano(),
			Fields: EventFields{
				BatteryLevel: int(rand.Int31n(101)),
			},
		}
		if bytes, err = json.Marshal(&ev); err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}
		if _, err = jc.PublishAsync(StreamName+"."+BatteryLevelSubject, bytes); err != nil {
			fmt.Println("Error: " + err.Error())
			continue
		}

		fmt.Printf("Published batch %d\n", i+1)
	}

	return nil
}

func GetLocationsSubscription(jc nats.JetStreamContext) (*nats.Subscription, error) {
	sub, err := jc.Subscribe(StreamName+"."+LocationSubject,
		handleLocations2,
		nats.MaxDeliver(3),
		nats.MaxAckPending(5),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create location subscription: %w", err)
	}

	return sub, nil
}

func handleLocations2(msg *nats.Msg) {
	var ev Event

	err := json.Unmarshal(msg.Data, &ev)
	if err != nil {
		fmt.Printf(">> Consumer: Location >> Error: %s\n", err.Error())
		msg.Nak()
		return
	}

	fmt.Printf(">> Consumer: Location >> Values:\n\tLongitude: %f\n\tLatitude: %f\n",
		ev.Fields.Longitude,
		ev.Fields.Latitude)

	msg.Ack()
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer nc.Close()

	jc, err := nc.JetStream()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	select {
	case <-jc.PublishAsyncComplete():
	}

	sub, err := jc.Subscribe(StreamName+"."+LocationSubject,
		handleLocations,
	)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	go gracefulShutdown(sub)
	forever := make(chan int)
	<-forever
}

func gracefulShutdown(sub *nats.Subscription) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		sub.Unsubscribe()
		sub.Drain()
		fmt.Println("Shutting down gracefully.")
		// clean up here
		os.Exit(0)
	}()
}

// LocalMilliscond 当前毫秒
func LocalMilliscond() int64 {
	return time.Now().UnixNano() / 1e6
}

func handleLocations(msg *nats.Msg) {
	//fmt.Printf(">> Consumer: Location >> Payload:%s\n", string(msg.Data))
	v := &Event{}
	if err := json.Unmarshal(msg.Data, v); err == nil {
		fmt.Printf("cost %d ms \n", (LocalMilliscond() - v.Timestamp))
	}
	msg.Ack()
}
