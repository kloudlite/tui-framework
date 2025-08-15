package main

import (
    "bufio"
    "fmt"
    "os"

    "github.com/kloudlite/tui-framework/pkg/state"
)

func main() {
    counter := state.NewReactive[int](0)

    // subscribe to counter changes
    counter.Subscribe(func(value int) {
        fmt.Printf("\rCounter: %d", value)
    })

    fmt.Println("Press Enter to increment the counter. Type 'q' and press Enter to exit.")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        if !scanner.Scan() {
            break
        }
        input := scanner.Text()
        if input == "q" {
            break
        }
        // increment the counter
        current := counter.Get()
        counter.Set(current + 1)
    }
}
