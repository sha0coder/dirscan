package main

import "os"
import "fmt"
import "bufio"

type Wordlist struct {
    Words []string
}

func (w *Wordlist) Clean() {
    // implement this ;)
}

func (w *Wordlist) Load(filename string) {
    file, err := os.Open(filename)
    check(err, "Can't load the wordlist")
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
         w.Words = append(w.Words, scanner.Text())
    }
    fmt.Println("Wordlist Loaded.")
}

