package main

import (
	"crypto/md5"
	"fmt"
	"time"
	"math/rand"
)

type request_type int
const(put_type request_type = 1
	  pull_type request_type = 2
	  kill_type request_type = 3
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
	msg			message_type
	key_values  []key_value
}

var kv_channels []chan key_value
var node_channels []chan message
var response_channels []chan key_value
var num_nodes int

func main() {
	num_nodes = 5
	//Initializing Channels
	fmt.Println("Initializing channels for communication")
	//These channels are to send key values
	for i := 0; i < num_nodes; i++ {
		temp_kv_chan := make(chan key_value, 100)
		kv_channels = append(kv_channels, temp_kv_chan)

		temp_node_chan := make(chan message, 100)
		node_channels = append(node_channels, temp_node_chan)

		temp_resp_chan := make(chan key_value, 100)
		response_channels = append(response_channels, temp_resp_chan)
	}

	for i := 0; i < num_nodes; i++ {
		go bucket(i)
	}

	var put_kv [5] key_value

	put_kv[0].key = "Maria"
	put_kv[0].value = 100
	put_kv[0].req = put_type

	put_kv[1].key = "John"
	put_kv[1].value = 20
	put_kv[1].req = put_type

	put_kv[2].key = "Anna"
	put_kv[2].value = 40
	put_kv[2].req = put_type

	put_kv[3].key = "Tim"
	put_kv[3].value = 100
	put_kv[3].req = put_type

	put_kv[4].key = "Alex"
	put_kv[4].value = 10
	put_kv[4].req = put_type

	var pull_kv [5] key_value

	pull_kv[0].key = "Maria"
	pull_kv[0].value = 0
	pull_kv[0].req = pull_type

	pull_kv[1].key = "John"
	pull_kv[1].value = 0
	pull_kv[1].req = pull_type

	pull_kv[2].key = "Anna"
	pull_kv[2].value = 0
	pull_kv[2].req = pull_type

	pull_kv[3].key = "Tim"
	pull_kv[3].value = 0
	pull_kv[3].req = pull_type

	pull_kv[4].key = "Alex"
	pull_kv[4].value = 0
	pull_kv[4].req = pull_type



	kill_chan_timeout := time.NewTimer(15 * time.Second)
	send_timeout := time.NewTimer(1 * time.Second)
	curr_kv_send := 0
	put_pull := 1
	for {
		select {
		case <-kill_chan_timeout.C:
			n := rand.Intn(num_nodes)
			var kv key_value
			kv.req = kill_type
			kv_channels[n] <- kv
			kill_chan_timeout.Reset(15 * time.Second)

		case <-send_timeout.C:
			sent := false
			hash := md5.Sum([]byte(put_kv[curr_kv_send].key))
			var n int = (((int(hash[1]) << 8) | int (hash[0])) % 360) / (360 / num_nodes)
			if put_pull == 1{
				fmt.Println("Sending Key: ", put_kv[curr_kv_send].key, " Value: ", put_kv[curr_kv_send].value, " to node ", n)
				kv_channels[n] <- put_kv[curr_kv_send]
			}
			if put_pull == 0{
				fmt.Println("Pulling Key: ", pull_kv[curr_kv_send].key)
				kv_channels[n] <- pull_kv[curr_kv_send]
			}


			resp_timeout := time.NewTimer(100 * time.Millisecond)
			for ; sent == false; {
				select {
				case <-resp_timeout.C:
					fmt.Println("Node: ", n, " did not respond")
					n = (n + 1) % num_nodes
					if put_pull == 1{
						kv_channels[n] <- put_kv[curr_kv_send]
					}else{
						kv_channels[n] <- pull_kv[curr_kv_send]
					}
					resp_timeout.Reset(100 * time.Millisecond)

				case resp := <- response_channels[n]:
					resp_timeout.Stop()
					sent = true
					if resp.req == pull_type {
						fmt.Println("Pull successful from node:", n, " Key: ", resp.key, " Value: ", resp.value)
					} else if resp.req == put_type && resp.key == put_kv[curr_kv_send].key && resp.value == put_kv[curr_kv_send].value{
						fmt.Println("Put successful into node ", n)
					}

				default:
					break
				}
			}
			send_timeout.Reset(1 * time.Second)
			curr_kv_send = curr_kv_send + 1
			if curr_kv_send == 5{
				curr_kv_send = 0
				if put_pull == 0{
					put_pull = 1
				}else{
					put_pull = 0
				}
			}

		default:
			break
		}

	}

}

func bucket(node_number int) {
	var bucket_values map[string]int = make(map[string]int)
	var return_kv key_value
	prev_node := node_number - 1
	if prev_node < 0 {
		prev_node = num_nodes - 1
	}
	fmt.Println("Launching node: ", node_number)
	for {
		select {
		case put_kv := <-kv_channels[node_number]:
			if put_kv.req == put_type {
				bucket_values[put_kv.key] = put_kv.value
				var kv key_value
				kv.req = put_type
				kv.value = put_kv.value
				kv.key = put_kv.key
				response_channels[node_number] <- kv
			} else if put_kv.req == pull_type {
				return_kv.req = pull_type
				return_kv.key = put_kv.key
				return_kv.value = bucket_values[put_kv.key]
				response_channels[node_number] <- return_kv
			} else if put_kv.req == kill_type {
				// send replicate request and sleep, then send update request
				var replicate_msg message
				replicate_msg.node_id = node_number
				replicate_msg.msg = replicate_type
				for k, v := range bucket_values {
					var kv key_value
					kv.key = k
					kv.value = v
					kv.req = put_type
					replicate_msg.key_values = append(replicate_msg.key_values, kv)
				}
				node_channels[node_number] <- replicate_msg
				fmt.Println("Node: ", node_number, " going to sleep")
				suspend()
				fmt.Println("Node: ", node_number, " waking up")
				bucket_values = make(map[string]int)
				done := false
				for ; done == false; {
					select {
					case <- kv_channels[node_number]:
						break
					default:
						done = true
						break
					}
				}
				var update_req message
				update_req.node_id = node_number
				update_req.msg = update_type
				node_channels[node_number] <- update_req
			}
		case input_msg := <-node_channels[prev_node]:
			if input_msg.msg == replicate_type {
				// receive all data from failing node and store in our own
				for _, val := range input_msg.key_values {
					bucket_values[val.key] = val.value
				}

			} else if input_msg.msg == update_type {
				// send all data back to failed node when back online
				var resp message
				resp.node_id = node_number
				resp.msg = response_type
				for k, v := range bucket_values {
					hash := md5.Sum([]byte(k))
					if (((int(hash[1]) << 8) | int(hash[0])) % 360) < ((input_msg.node_id * (360/num_nodes)) - 1) {
						var kv key_value
						kv.key = k
						kv.value = v
						kv.req = put_type
						resp.key_values = append(resp.key_values, kv)
					}
				}
				node_channels[prev_node] <- resp
			}

		case input_msg := <-node_channels[node_number]:
			if input_msg.msg == response_type {
				for _, val := range input_msg.key_values {
					bucket_values[val.key] = val.value
				}
			}
		}
	}
}

func suspend() {
	suspend_timeout := time.NewTimer(10 * time.Second)
	for {
		select {
		case <- suspend_timeout.C:
			return

		default:
			break
		}
	}
}
