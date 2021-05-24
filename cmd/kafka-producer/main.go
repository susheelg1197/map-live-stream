package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
type Property struct {
}
type Feature struct {
	Type       string   `json:"type"`
	Properties Property `json:"properties"`
	Geometry   Geometry `json:"geometry"`
}
type GeoJson struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var jsonString = `{
	"type": "FeatureCollection",
	"features": [
	  {
		"type": "Feature",
		"properties": {},
		"geometry": {
		  "type": "Polygon",
		  "coordinates": [
			
			  [
				73.63037109375,
				18.578568865536027
			  ],
			  [
				73.76220703125,
				18.445741249840815
			  ],
			  [
				73.86657714843749,
				18.323240460443397
			  ],
			  [
				74.13711547851562,
				18.294557510034192
			  ],
			  [
				74.29504394531249,
				18.394927021680232
			  ],
			  [
				74.32388305664062,
				18.497842796761336
			  ],
			  [
				74.3994140625,
				18.603299856818385
			  ],
			  [
				74.37606811523438,
				18.720397774505507
			  ],
			  [
				74.29641723632812,
				18.802318121688117
			  ],
			  [
				74.20852661132812,
				18.846512551704713
			  ],
			  [
				74.168701171875,
				18.903688072314985
			  ],
			  [
				74.036865234375,
				18.95954527856256
			  ],
			  [
				73.86795043945312,
				18.94785578129413
			  ],
			  [
				73.65234375,
				18.929670491513274
			  ],
			  [
				73.52325439453124,
				18.80751806940863
			  ],
			  [
				73.49441528320312,
				18.734704192913256
			  ],
			  [
				73.45046997070312,
				18.626725903064777
			  ],
			  [
				73.44635009765625,
				18.470491457966812
			  ],
			  [
				73.49166870117188,
				18.376682358161855
			  ],
			  [
				73.54110717773438,
				18.2867340628812
			  ],
			  [
				73.63311767578125,
				18.215002696822772
			  ],
			  [
				73.75808715820312,
				18.203262019544418
			  ],
			  [
				73.85421752929688,
				18.204566578321632
			  ],
			  [
				74.02587890625,
				18.195434461776465
			  ],
			  [
				73.63037109375,
				18.578568865536027
			  ]
			
		  ]
		}
	  }
	]
  }`

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
					fmt.Println("Delivered message:: ", ev.TopicPartition, string(ev.Value))
				}
			}
		}
	}()

	topic := "map_stream"
	var resp GeoJson
	json.Unmarshal([]byte(jsonString), &resp)
	// fmt.Println(resp.Features[0].Geometry.Coordinates[1])
	for _, coordinate := range resp.Features[0].Geometry.Coordinates {

		jsonString, _ := json.Marshal(map[string]interface{}{
			"longitude": coordinate[0],
			"latitude":  coordinate[1],
		})
		// fmt.Println(string(jsonString))

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: jsonString,
		}, nil)
		if err != nil {

			fmt.Println(err)
		}
		p.Flush(15 * 1000)
	}

	// for i := 0; i < 10; i++ {
	// 	jsonStructure, _ := json.Marshal(map[string]interface{}{
	// 		"name":    "Susheel",
	// 		"message": "hi",
	// 		"number":  100,
	// 	})
	// 	var b []byte
	// 	b = append(b, byte(i))
	// 	err := p.Produce(&kafka.Message{
	// 		TopicPartition: kafka.TopicPartition{
	// 			Topic:     &topic,
	// 			Partition: kafka.PartitionAny,
	// 		},
	// 		Value: jsonStructure,
	// 	}, nil)
	// 	if err != nil {

	// 		fmt.Println(err)
	// 	}
	// 	p.Flush(15 * 1000)
	// }
}
