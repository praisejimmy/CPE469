package main

import ("fmt"
        "time"
        "bufio"
        "os"
        "strings"
        "runtime"
        "sync"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Format: ./word_search.exe input_file search_string")
        os.Exit(-1)
    }
    file_name := os.Args[1]
    str := os.Args[2]
    fmt.Println("Testing sequential string search")
    t := time.Now()
    fmt.Println("Matching words found: ", search_string_seq(file_name, str))
    fmt.Println(time.Since(t))
    fmt.Println("Testing concurrent string search")
    t = time.Now()
    fmt.Println("Matching words found: ", search_string_conc(file_name, str))
    fmt.Println(time.Since(t))
}

func search_string_seq(file_name string, str string) int{
    str_freq := 0
    str_len := len(str)
    f, err := os.Open(os.Args[1])
    if err != nil {
        fmt.Println("Unable to open file")
        os.Exit(-1)
    }
    defer f.Close()
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        txt_line := scanner.Text()
        words := strings.Split(txt_line, " ")
        for i := 0; i < len(words); i++ {
            if strncmp(str, words[i], str_len) {
                str_freq++
            }
        }
    }
    return str_freq
}

func strncmp(str1 string, str2 string, bytes int) bool {
    len1 := len(str1)
    len2 := len(str2)
    for i := 0; i < bytes; i++ {
        if i >= len1 || i >= len2 || str1[i] != str2[i] {
            return false
        }
    }
    return true
}

func search_string_conc(file_name string, str string) int {
    input_ch := make(chan []string, 8)
    freq_ch := make(chan int, 8)
    var wg sync.WaitGroup
    str_freq := 0
    f, err := os.Open(os.Args[1])
    if err != nil {
        fmt.Println("Unable to open file")
        os.Exit(-1)
    }
    defer f.Close()
    for i := 0; i < runtime.NumCPU(); i++ {
        go func(input_ch chan []string, freq_ch chan int) {
            wg.Add(1)
            for text := range input_ch {
                partial_count := 0
                for j := 0; j < len(text); j++ {
                    if strncmp(str, text[j], len(str)) {
                        partial_count++
                    }
                }
                freq_ch <- partial_count
            }
            wg.Done()
        }(input_ch, freq_ch)
    }
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        txt_line := scanner.Text()
        words := strings.Split(txt_line, " ")
        input_ch <- words
        str_freq += <-freq_ch
    }
    close(input_ch)
    wg.Wait()
    close(freq_ch)
    return str_freq
}
