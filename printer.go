
package main 

import "fmt"


type Printer struct {

}

/*
type Color struct {
    Clean: "\\033[0m",
    Clear: "\\33[2K",
    Bold:  "\\033[22m",
    Italic: "\\033[23m']",
    Underline: "\\033[24m']",
    Inverse: "\\033[27m']",
    White: "\\033[39m']",
    Grey: "\\033[39m']",
    Black: "\\033[39m']",
    Blue: "\\033[39m']",
    Cyan: "\\033[39m']",
    Green: "\\033[39m']",
    Magenta: "\\033[39m']",
    Red: "\\033[39m']",
    Yellow: "\\033[39m"
}*/

func (p *Printer) Show(b string, code int, size int, url string) {
    
    if code == 200 {
        fmt.Printf("%s:[%d] (%d bytes)\t%s\n",b,code,size,url)
    
    } else if 301 <= code && code <= 303 {
        fmt.Printf("%s:[%d] (redirect)\t%s\n",b,code,url) //TODO:where?

    } else if code == 401 {
        fmt.Printf("%s:[%d] (auth needed)\t%s\n",b,code,url)

    } else if code == 403 {
        fmt.Printf("%s:[%d] (denied)\t%s\n",b,code,url)

    }

}

