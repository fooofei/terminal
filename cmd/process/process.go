package main

import (
	"fmt"
	terminal "github.com/fooofei/terminal/writer"
	"time"
)

func main() {
	writer, err := terminal.New()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i = i + 10 {
		// add your text to writer's buffer
		writer.Write(fmt.Appendf(nil, "Downloading (%d/100) bytes...", i))
		time.Sleep(time.Millisecond * 200)

		// clear the text written by previous write, so that it can be re-written.
		writer.Clear()
	}

	// reset the writer
	writer.Clear()
	_ = writer.Close()
	fmt.Println("Download finished!")

}
