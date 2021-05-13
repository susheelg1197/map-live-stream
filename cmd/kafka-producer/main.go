package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:29092"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Println("Delivery  Failed", ev.TopicPartition)
				} else {
					fmt.Println("Delivered message:: ", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "map_stream"
	for i := 0; i < 10; i++ {
		var b []byte
		b = append(b, byte(i))
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: b,
		}, nil)
		if err != nil {

			fmt.Println(err)
		}
		p.Flush(15 * 1000)
	}
}
