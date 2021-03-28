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
	fmt.Fprintf(w, "hello1\n")
	fmt.Fprintf(w, "hello2\n")
	fmt.Fprintf(w, "hello3\n")
	fmt.Fprintf(w, "hello4\nhello44 ")
	fmt.Fprintf(w, "hello5 ")
	fmt.Fprintf(w, "hello6 ")
	fmt.Fprintf(w, "hello7\n")

	time.Sleep(time.Second)
	_ = w.Close()

}
