package main

import ("fmt"
        "time"
        "math/rand"
)

type member struct {
    id int
    hb int
    ts time.Time
}

var chans [8]chan []member
var kill_chans [8]chan bool

func main() {
    fmt.Println("Beginning membership protocol")
    for i := range chans {
        chans[i] = make(chan []member, 4)
    }
    for i := range kill_chans {
        kill_chans[i] = make(chan bool)
    }
    for i := 0; i < len(chans); i++ {
        go member_routine(i)
    }
    kill_timeout := time.NewTimer(20 * time.Second)
    for{
        select {
        case <-kill_timeout.C:
            rand_id := rand.Intn(3)
            kill_chans[rand_id] <- true
            kill_timeout = time.NewTimer(20 * time.Second)

        default:

        }
    }
}

func member_routine(my_id int) {
    delete_chan := make(chan bool, 4)
    var membership_table []member
    var delete_list []int
    table_chan1 := my_id
    var table_chan2 int
    if (table_chan1 == 0) {
        table_chan2 = 7
    } else {
        table_chan2 = table_chan1 - 1
    }
    my_membership := member {
        id: my_id,
        hb: 0,
        ts: time.Now(),
    }
    membership_table = append(membership_table, my_membership)
    hb_timeout := time.NewTimer(10 * time.Millisecond)
    share_timeout := time.NewTimer(20 * time.Millisecond)
    for {
        select {
        case <-hb_timeout.C:
            for i := range membership_table {
                if membership_table[i].id == my_id {
                    membership_table[i].hb++
                    membership_table[i].ts = time.Now()
                }
            }
            hb_timeout = time.NewTimer(10 * time.Millisecond)

        case <-share_timeout.C:
            if rand.Intn(2) == 0 {
                chans[table_chan1] <- membership_table
            } else {
                chans[table_chan2] <- membership_table
            }
            share_timeout = time.NewTimer(20 * time.Millisecond)

        case other_table := <-chans[table_chan1]:
            var to_delete []int
            membership_table, to_delete = merge_tables(my_id, membership_table, other_table)

            for i := range to_delete {
                if !contains(delete_list, to_delete[i]) {
                    delete_list = append(delete_list, to_delete[i])
                    go func() {
                        delete_timeout := time.NewTimer(10 * time.Second)
                        <-delete_timeout.C
                        delete_chan <- true
                    }()
                }
            }

        case other_table := <-chans[table_chan2]:
            var to_delete []int
            membership_table, to_delete = merge_tables(my_id, membership_table, other_table)

            for i := range to_delete {
                if !contains(delete_list, to_delete[i]) {
                    delete_list = append(delete_list, to_delete[i])
                    go func() {
                        delete_timeout := time.NewTimer(5 * time.Second)
                        <-delete_timeout.C
                        delete_chan <- true
                    }()
                }
            }

        case <-delete_chan:
            fmt.Printf("\nMember ID %d seen as offline by member ID %d\n", delete_list[0], my_id)
            membership_table = clean_member_list(membership_table, delete_list[0])
            fmt.Println("Member list after removing")
            print_table(membership_table)
            fmt.Println()
            delete_list = delete_list[1:]

        case <-kill_chans[my_id]:
            fmt.Printf("********************************************Killing member ID: %d\n", my_id)
            time.Sleep(15 * time.Second)
            membership_table = []member{}
            membership_table = append(membership_table, my_membership)
            fmt.Printf("********************************************Member online ID: %d\n", my_id)

        default:

        }
    }
}

func print_table(table []member) {
    for i := range table {
        fmt.Printf("\tID: %d, HB: %d, Time: %d\n", table[i].id, table[i].hb, table[i].ts.Unix())
    }
}

func clean_member_list(my_table []member, delete_node int) []member {
    for i := range my_table {
        if delete_node == my_table[i].id {
            my_table[i] = my_table[len(my_table) - 1]
            my_table = my_table[:len(my_table) - 1]
            break
        }
    }
    return my_table
}

func merge_tables(my_id int, my_table []member, other_table []member) ([]member, []int) {
    var to_delete []int
    for i := range other_table {
        found := false
        for j := range my_table {
            if (my_table[j].id == other_table[i].id) {
                found = true
            }
            if found && my_table[j].hb < other_table[i].hb {
                my_table[j].hb = other_table[i].hb
                my_table[j].ts = time.Now()
            } else if found && (my_table[j].hb == other_table[i].hb) && (time.Since(my_table[j].ts) > (3 * time.Second)) {
                to_delete = append(to_delete, my_table[j].id)
            }
            if found {
                break
            }
        }
        if !found {
            fmt.Printf("Member ID %d seen as online by member ID %d\n", other_table[i].id, my_id)
            my_table = append(my_table, other_table[i])
        }
    }
    return my_table, to_delete
}

func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
