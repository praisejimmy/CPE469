package main

import (
	//"crypto/md5"
	"fmt"
	//"strings"
)

type request_type int
const(put_type request_type = 1
	  pull_type request_type = 2
)

type key_value struct {
	key          string
	value        int
	req 		request_type
}

type message_type int
const(update_type message_type = 1
	  replicate_type message_type = 2
	  response_type message_type = 3
)

type message struct {
	node_id		int
	mes			message_type
	key_values  []key_value
}

func main() {
	num_nodes := 5

	//Initializing Channels
	fmt.Println("Initializing channels for communication")
	//These channels are to send key values
	var kv_channels []chan key_value
	for i := 0; i < num_nodes; i++ {
		temp_chan := make(chan key_value)
		kv_channels = append(kv_channels, temp_chan)
	}

	for i := 0; i < num_nodes; i++ {
		go bucket(kv_channels[i], (i+1) * 71)
	}

	var input_kv key_value

	input_kv.key = "Martha"
	input_kv.value = 10
	input_kv.req = put_type
	kv_channels[0] <- input_kv

	input_kv.key = "Martha"
	input_kv.value = 0
	input_kv.req = pull_type

	kv_channels[0] <- input_kv

	for response := range kv_channels[0] {
		fmt.Println("Key: ", response.key, " Value: ", response.value)
	}

	for {
	}

}

func bucket(my_kv_chan chan key_value, node_number int) {
	var bucket_values map[string]int = make(map[string]int)
	var return_kv key_value
	fmt.Println("Lanching bucket: ", node_number) 
	for {
		select {
		case input_kv := <-my_kv_chan:
			if input_kv.req == put_type {
				bucket_values[input_kv.key] = input_kv.value
			} else if input_kv.req == pull_type {
				return_kv.req = pull_type
				return_kv.key = input_kv.key
				return_kv.value = bucket_values[input_kv.key]
				my_kv_chan <- return_kv
			}
		}
	}
}
