package main

import (
	"fmt"
	terminal "github.com/fooofei/terminal/writer"
	"time"
)

func main() {
	w, err := terminal.New()
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintf(w, "hello1\n")
	_, _ = fmt.Fprintf(w, "hello2\n")
	_, _ = fmt.Fprintf(w, "hello3\n")
	_, _ = fmt.Fprintf(w, "hello4\nhello44 ")
	_, _ = fmt.Fprintf(w, "hello5 ")
	_, _ = fmt.Fprintf(w, "hello6 ")
	_, _ = fmt.Fprintf(w, "hello7\n")

	time.Sleep(time.Second)
	_ = w.Close()

}
