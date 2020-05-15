package main

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
)

func main(){
    var wg sync.WaitGroup
    var chans [5]chan [5]Heartbeat_Table
    fail_channel := make(chan int)
    //Initialize channels to publish tables on
    for i := range chans {
        chans[i] = make(chan [8] Heartbeat_Table, 8)
    }
    wg.Add(1)
    //Go routine that sends a message to fail on a channel
    //Random because don't know which node will read first.
    go func(fail_channel chan int){
        fail_timer := time.NewTimer(15 * time.Second)
        fmt.Println("Going to start kill timer")
        for{
            select{
                case <- fail_timer.C:
                    fmt.Println("Sending fail notice")
                    fail_channel <- 1
                    fail_timer = time.NewTimer(15 * time.Second)
            }
        }
        wg.Done()
    }(fail_channel)
    //Launch 8 computing nodes with the correct neighbor channels.
    for i := 0; i < 8; i ++{
        table_temp := Heartbeat_Table{i, 1, time.Now()}
        if i == 0{
            wg.Add(1)
            go node(chans[i], table_temp, &wg, chans[7], chans[1], fail_channel)
        }else if(i == 7){
            wg.Add(1)
            go node(chans[i], table_temp, &wg, chans[6], chans[0], fail_channel)
        }else{
            wg.Add(1)
            go node(chans[i], table_temp, &wg, chans[i-1], chans[i+1], fail_channel)
        }
    }
    wg.Wait()

    fmt.Println("Finished main")
}
type key_value struct {
	key string
  value int
  request_type int
}

func main() {
	num_nodes := 5

	//Initializing Channels
	fmt.Println("Initializing channels for communication")
	//These channels are to send key values
	var kv_channels [num_nodes]chan key_value
	for i := range kv_channels {
		kv_channels[i] = make(chan key_value)
	}

}

func bucket(my_kv_chan chan key_value, node_number int){
  put := 1
  pull := 2
  var bucket_values map[string]int
  for {
    select{
    case input_kv := <- my_kv_chan:
      if input_kv.request_type == put{

      }
    }
  }
}
