/*
    dirscan fast web resource enumerator
    crawler + bruterforcer
    @sha0coder

    TODO:
    - store log
    - get redirection to brute/crawl/subparts
    - colors

*/

package main

import "os"
import "fmt"
import "flag"
import "errors"
import "strings"


var R *Requests
var P *Printer
var C *Crawl
var W *Wordlist
var CFG Config
var el_c int = 0

var ext_jsp = []string{"html","jsp","do","cfg","prop","sql","log","txt","zip","rar","tar","7z","mdb","tar.gz","tar.bz2","pem","class","jar"} 
var ext_asp = []string{"html","asp","aspx","mdb","cfg","conf","ini","log","txt","log","zip","rar","7z","mdb","pem"}
var ext_php = []string{"html","php","php4","sql","cfg","conf","log","txt","zip","rar","tar","7z","mdb","tar.gz","tar.bz2","pem"}


func check(err error, msg string) {
    if err != nil {
        fmt.Println(msg)
        os.Exit(1)
    }
}

func checkWebserver(surl string) {
    _, code, resp := R.Get(surl)
    R.QuitOnFail(code, "Can't connect")
    fmt.Printf("Server: %s\nDefault response: %d\n", resp.Header.Get("Server"), resp.StatusCode)

    _, code, resp = R.Options(surl)
    fmt.Println(code);
    if (code > 0) {
       fmt.Printf("Allowed Options: %s\n", resp.Header.Get("Allow"))
    }
}

func IsDirectory(url string) bool {
    parts := strings.Split(url, "/")
    l := len(parts)
    if l>2 {
        if strings.Contains(parts[l-1], ".") {
            return false
        }
    }
    return true
}

func EndLogic(res *[]string) {
    el_c++
    //fmt.Printf("EL start %d\n",el_c) // recursion debug
    for _, u := range *res {
        if IsDirectory(u) {
            //brute
            b := new(Bruter)
            b.OnEnd(EndLogic)
            b.Brute(u+"/")

        } else {
            //crawl
            C.AddHost(CFG.Host)
            C.Crawler(u)
            if len(C.NewResources)>0 {
                EndLogic(&C.NewResources)
            }
        }
    }
    //fmt.Printf("EL end %d\n",el_c) // recursion debug
    el_c--
    if el_c <= 0 { 
        fmt.Println("done.")
        //os.Exit(1) 
    }
}

func main() {
    var i int
    var url *string = flag.String("url", "", "the url")
    var wordlist *string = flag.String("dict", "", "the wordlist")
    var goroutines *int = flag.Int("go", 5, "num of concurrent goroutines")
    var platform *string = flag.String("lang", "", "languaje (java, asp or php)")
    var proxy *string = flag.String("proxy", "", "set proxy ip:port")
    flag.Parse()

    if *url == "" || *wordlist == "" || *platform == "" {
        check(errors.New(""), "bad usage,  --help")
    }

    fmt.Printf("url:[%s]\n", *url)

    switch (*platform) {
        case "java":
            CFG.Extensions = ext_jsp
        case "asp":
            CFG.Extensions = ext_jsp
        case "php":
            CFG.Extensions = ext_jsp
        default:
            check(errors.New(""), "bad platform,  --help")
    }

    R = NewRequests()
    C = NewCrawl()
    P = new(Printer)
    W = new(Wordlist)
    
    CFG.Goroutines = *goroutines 
    CFG.Url = *url 
    CFG.Host = strings.Split(*url,"/")[2]

    if *proxy != "" { R.SetProxy("http://" + *proxy) }
    checkWebserver(*url)
    if !strings.HasSuffix(*url,"/") { *url = *url+"/" }

    W.Load(*wordlist)

    B := new(Bruter)
    B.OnEnd(EndLogic)
    B.Brute(CFG.Url)

    fmt.Printf("Scanning, press enter to interrupt.\n")
    fmt.Scanf("%d", &i)
    fmt.Printf("interrupted.")

}
