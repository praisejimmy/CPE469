go func(txt_line []byte, ch chan int) {
    partial_freq := 0
    curr_str := ""
    line_len := len(txt_line)
    for i := 0; i < line_len; i++ {
        if unicode.IsSpace(rune(txt_line[i])) || unicode.IsPunct(rune(txt_line[i])) {
            if curr_str == str {
                fmt.Println("Curr: ", curr_str, " Str: ", str)
                partial_freq++
            }
            curr_str = ""
        } else {
            curr_str = curr_str + string(txt_line[i])
        }
    }
    ch<-partial_freq
}(txt_curr, ch)
