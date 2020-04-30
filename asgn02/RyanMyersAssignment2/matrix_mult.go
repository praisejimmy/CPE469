package main
import ("fmt"
        "math/rand"
        "time"
        "sync"
        "runtime"
)

func main() {
    var mtx1 [1000][1000]float64
    var mtx2 [1000][1000]float64
    var i, j int
    rand.Seed(time.Now().UnixNano())
    for i = 0; i < 1000; i++ {
        for j = 0; j < 1000; j++ {
            mtx1[i][j] = rand.Float64()
            mtx2[i][j] = rand.Float64()
        }
    }
    fmt.Println("Timing sequential array multiplication")
    start := time.Now()
    mult_seq(mtx1, mtx2)
    fmt.Println(time.Since(start))

    fmt.Println("Timing concurrent array multiplication")
    start = time.Now()
    mult_conc(mtx1, mtx2)
    fmt.Println(time.Since(start))
}

func mult_seq(mtx1 [1000][1000]float64, mtx2 [1000][1000]float64) [1000][1000]float64{
    var i, j, k int
    var result [1000][1000]float64
    for i = 0; i < 1000; i++ {
        for j = 0; j < 1000; j++ {
            for k = 0; k < 1000; k++ {
                result[i][j] += mtx1[i][k] * mtx2[k][j]
            }
        }
    }
    return result
}

func mult_conc(mtx1 [1000][1000]float64, mtx2 [1000][1000]float64) [][]float64{
    var wg sync.WaitGroup
    max_procs := runtime.NumCPU()
    ret := make([][]float64, 1000)
    for i := range ret {
        ret[i] = make([]float64, 1000)
    }
    for i := 0; i < max_procs; i++ {
        wg.Add(1)
        if i == (max_procs - 1) {
            go func(ret [][]float64, start int, rows int) {
                defer wg.Done()
                for i := start; i < start + rows; i++ {
                    for j := 0; j < 1000; j++ {
                        for k := 0; k < 1000; k++ {
                            ret[i][j] += mtx1[i][k] * mtx2[k][j]
                        }
                    }
                }
            }(ret, i * (1000/max_procs), 1000 - i * (1000/max_procs))
        } else {
            go func(ret [][]float64, start int, rows int) {
                defer wg.Done()
                for i := start; i < start + rows; i++ {
                    for j := 0; j < 1000; j++ {
                        for k := 0; k < 1000; k++ {
                            ret[i][j] += mtx1[i][k] * mtx2[k][j]
                        }
                    }
                }
            }(ret, i * (1000/max_procs), 1000/max_procs)
        }
    }
    wg.Wait()
    return ret
}
