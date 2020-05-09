package main
import ("fmt"
        "runtime"
        "sync"
        "math/rand"
        "sort"
)

func main() {

    // mtx1 := make([][]int, 1000)
    // for i := range mtx1 {
    //     mtx1[i] = make([]int, 1000)
    // }
    //
    // var i, j int
    // for i = 0; i < 1000; i++ {
    //     for j = 0; j < 1000; j++ {
    //         mtx1[i][j] = rand.Int()
    //     }
    // }
    // find_max(mtx1)
    var i int
    mtx2 := make([]int, 100)
    for i = 0; i < 100; i++ {
        mtx2[i] = rand.Int()
    }
    fmt.Println(find_median(mtx2))
}

func find_max(Matrix [][]int) int {
    var wg sync.WaitGroup
    num_rows := len(Matrix)
    num_routines := runtime.NumCPU()
    max_chan := make(chan int, num_routines)
    for i := 0; i < num_routines; i++ {
        wg.Add(1)
        var curr_mtx [][]int
        if i == num_routines - 1 {
            curr_mtx = Matrix[i * (num_rows / num_routines):][:]
        } else {
            curr_mtx = Matrix[i * (num_rows / num_routines):(i + 1) * (num_rows / num_routines)][:]
        }
        go func (mtx [][]int, c chan int) {
            defer wg.Done()
            max := mtx[0][0]
            for i := 0; i < len(mtx); i++ {
                for j := 0; j < len(mtx[0]); j++ {
                    if mtx[i][j] > max {
                        max = mtx[i][j]
                    }
                }
            }
            c <- max
        } (curr_mtx, max_chan)
    }
    wg.Wait()
    global_max := <- max_chan
    done := false
    for ; done == false; {
        select {
            case curr_max := <- max_chan:
                if curr_max > global_max {
                    global_max = curr_max
                }

            default:
                done = true
                break
        }
    }
    return global_max
}

func find_median(Matrix []int) int {
    var wg sync.WaitGroup
    array_size := len(Matrix)
    num_routines := runtime.NumCPU()
    median_chan := make(chan int, num_routines)
    for i := 0; i < num_routines; i++ {
        wg.Add(1)
        var partition []int
        if i == num_routines - 1 {
            partition = Matrix[i * (array_size / num_routines):]
        } else {
            partition = Matrix[i * (array_size / num_routines):(i + 1) * (array_size / num_routines)]
        }
        go func (sub_mtx []int, c chan int) {
            defer wg.Done()
            sort.Ints(sub_mtx)
            c <- sub_mtx[len(sub_mtx) / 2]
        } (partition, median_chan)
    }
    wg.Wait()
    var median_vals []int
    done := false
    for ; done == false; {
        select {
        case new_val := <- median_chan:
            median_vals = append(median_vals, new_val)

        default:
            done = true
            break
        }
    }
    sort.Ints(median_vals)
    return median_vals[len(median_vals) / 2]
}
