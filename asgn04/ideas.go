package main
import ("fmt"
        "math/rand"
        "time"
        "sync"
)

type State int
const (
    Follower State = 0
    Candidate State = 1
    Leader State = 2
)

type Action int
const (
    Leader_Ping Action = 0
    Vote_Request Action = 1
    Vote_Cast Action = 2
    Complete_Term Action = 3
)

type Message struct {
    sender_id int
    action Action
}

var curr_leader int
var curr_leader_mux sync.Mutex
var num_nodes int

var action_chans []chan Message
var resp_chans []chan Message

func main() {

    curr_leader = -1
    num_nodes = 8

    action_chans = make([]chan Message, 8)
    resp_chans = make([]chan Message, 8)

    for i := 0; i < num_nodes; i++ {
        action_chans[i] = make(chan Message, 8)
        resp_chans[i] = make(chan Message, 8)
    }

    for i := 0; i < num_nodes; i++ {
        go Node_Machine(i)
    }

    // Run protocol indefinitely
    for{}
}

func Node_Machine(my_id int) {
    curr_state := Follower
    rand.Seed(time.Now().UnixNano())

    for {
        if curr_state == Follower {
            curr_state = Follower_State(my_id)
        } else if curr_state == Candidate {
            curr_state = Candidate_State(my_id)
        } else if curr_state == Leader {
            curr_state = Leader_State(my_id)
        }
    }
}

func Follower_State(my_id int) State {
    ms_time := time.Duration(rand.Intn(151) + 150)
    leader_ping_timeout := time.NewTimer(ms_time * time.Millisecond)
    for {
        select {
        case msg := <- action_chans[my_id]:
            if msg.sender_id == curr_leader && msg.action == Leader_Ping {
                ms_time = time.Duration(rand.Intn(151) + 150)
                leader_ping_timeout = time.NewTimer(ms_time * time.Millisecond)
            } else if msg.action == Vote_Request {
                Vote(my_id, msg.sender_id)
                return Follower
            }

        case <- leader_ping_timeout.C:
            return Candidate

        default:
            break
        }
    }
}

func Vote(my_id int, requester int) {
    var vote_resp Message
    vote_resp.sender_id = my_id
    vote_resp.action = Vote_Cast
    resp_chans[requester] <- vote_resp
    fmt.Println("Node ", my_id, " voting for node ", requester)
    for {
        select {
        case resp := <- action_chans[my_id]:
            if resp.sender_id == requester && resp.action == Complete_Term {
                return
            }

        default:
            break
        }
    }
}

func Candidate_State(my_id int) State {
    fmt.Println("Node ", my_id, " up for candidacy")
    var vote_req Message
    vote_cnt := 1
    vote_req.sender_id = my_id
    vote_req.action = Vote_Request
    for i := 0; i < num_nodes; i++ {
        if i != my_id {
            select {
            case action_chans[i] <- vote_req:

            default:

            }
        }
    }

    term_timeout := time.NewTimer(40 * time.Millisecond)

    for {
        select {
        case response := <- resp_chans[my_id]:
            if response.action == Vote_Cast {
                vote_cnt++
            }

        case <- term_timeout.C:
            var end_term_message Message
            end_term_message.sender_id = my_id
            end_term_message.action = Complete_Term
            for j := 0; j < num_nodes; j++ {
                if j != my_id {
                    action_chans[j] <- end_term_message
                }
            }
            if vote_cnt > num_nodes / 2 {
                curr_leader_mux.Lock()
                curr_leader = my_id
                curr_leader_mux.Unlock()
                fmt.Println("New leader elected, ID: ", my_id)
                return Leader
            }
            return Follower

        default:
            break
        }
    }
}

func Leader_State(my_id int) State {
    ping_timeout := time.NewTimer(5 * time.Millisecond)
    drop_leadership_timeout := time.NewTimer(15 * time.Second)
    var ping_message Message
    ping_message.sender_id = my_id
    ping_message.action = Leader_Ping
    for {
        select {
        case <- ping_timeout.C:
            for i := 0; i < num_nodes; i++ {
                if i != my_id {
                    action_chans[i] <- ping_message
                }
            }
            ping_timeout = time.NewTimer(5 * time.Millisecond)

        case <- drop_leadership_timeout.C:
            fmt.Println("ID: ", my_id, " dropping leadership")
            return Follower

        default:
            break
        }
    }
}
