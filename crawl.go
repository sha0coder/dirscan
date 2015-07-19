package main

import "strings"
import "fmt"
import "net/url"
import "code.google.com/p/go.net/html"

type CrawlEndCallback func(c *Crawl)

type Crawl struct {
    BaseUrl string
	R       *Requests
	Crawled []string
    Resources []string
    NewResources []string
	Hosts   []string
    Host string
	Pending chan string
	DoStop  bool
    EndCB CrawlEndCallback
}

func NewCrawl() *Crawl {
	c := new(Crawl)
	c.R = NewRequests()
	c.Pending = make(chan string)
	c.DoStop = false
	return c
}

func (b *Crawl) IsRepeated(url string) bool {
    for _, u := range b.Resources {
        if u == url {
            return true
        }
    }
    return false
}

func (c *Crawl) OnEnd(callback CrawlEndCallback) {
    c.EndCB = callback
}

func (c *Crawl) AddHost(host string) {
	c.Hosts = append(c.Hosts, host)
}

func (c *Crawl) IsCrawled(url string) bool {
	for _, u := range c.Crawled {
		if url == u {
			return true
		}
	}
	return false
}

func (c *Crawl) IsAllowed(host string) bool {
	for _, h := range c.Hosts {
		if host == h {
			return true
		}
	}
	return false
}

func (c *Crawl) Queue(url string) {
	host := strings.Split(url, "/")
	if len(host) < 3 {
		//fmt.Println("bad url: " + url)
		return
	}

	if !c.IsAllowed(host[2]) {
		//fmt.Printf("out scope: " + url)
		return
	}

	if !c.IsCrawled(url) {
		//fmt.Println("queued: " + url)
		c.Pending <- url // crash 
	}
}

func (c *Crawl) Stop() {
	c.DoStop = false
}

func (c *Crawl) Scan(surl string) {
    //fmt.Printf("scanning %s\n",surl)

	resp := c.R.LaunchNoRead("GET", surl, "")
	if resp == nil || resp.Body == nil { 
        fmt.Println("nil response: "+surl)
        return
    }
    defer resp.Body.Close()

	page := html.NewTokenizer(resp.Body)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			c.Crawled = append(c.Crawled, surl)
			return
		}
		token := page.Token()

		//if tokenType == html.StartTagToken { //&& token.DataAtom.String() == "a" {
		for _, attr := range token.Attr {
			if attr.Key == "href" || attr.Key == "action" || attr.Key == "src" {
                res := c.FixUrl(attr.Val)
                if res != "" && !c.IsRepeated(res) {

                    oUrl, err := url.Parse(res)
                    if err == nil {

                        if oUrl.Host == c.Host {

                            var test string

                            idx := strings.LastIndex(oUrl.Path, ".")
                            if idx >=0 {
                                oUrl.Path = oUrl.Path[0:idx]+"test1337"+oUrl.Path[idx+1:]     //TODO: si la url acaba en punto, crashea out of index
                                test = oUrl.String()
                            } else {
                                test = res
                            }

                            //fmt.Printf("test:%s\n",test)
                            _, code_not_found, _ := R.Get(test)
                            html, code, _ := R.Get(res)

                            if code != code_not_found {
                                P.Show("c",code, len(html), res)
                                c.Resources = append(c.Resources, res)
                                c.NewResources = append(c.NewResources, res)
                            }

                        }
                    }
                }
			}
		}
	}
}

func (c *Crawl) FixUrl(r string) string {
    var durl string
    parts := strings.Split(c.BaseUrl,"/")
    dir := strings.Join(parts[0:len(parts)-1],"/")+"/" //posible crash

    if strings.HasPrefix(r,"../") {
        var nurl []string
        durl = dir+r //deberia pasrar de estos casos, tambien puede estar .. por medio
    
        for _,s := range strings.Split(durl,"/") {
            if s == ".." {
                nurl = nurl[0:len(nurl)-1]
            } else {
                nurl = append(nurl, s)
            } 
        }

        if len(nurl)<4 { return "" }

        durl = strings.Join(nurl, "/")

    } else if strings.HasPrefix(r,"/") {
        durl = dir+r

    } else if strings.HasPrefix(r,"./") {
        durl  = dir+r // quitar el punto

    } else if strings.HasPrefix(r,"http") {
        durl = r

    } else {
        durl = dir+r
    }


    oUrl, err := url.Parse(durl)
    if err != nil { 
        fmt.Println("bogus url: "+durl)
        return "" 
    }

    return oUrl.String()
}

func (c *Crawl) SubUrls() {
    for _, u := range c.Resources {
        var durl string
        spl := strings.Split(u,"/")
        sz := len(spl)
        for i:=4; i<sz; i++ {
            durl = strings.Join(spl[0:i],"/")
            if !c.IsRepeated(durl) {
                c.Resources = append(c.Resources, durl)
                c.NewResources = append(c.NewResources, durl)
            }
        }
    }
}

func (c *Crawl) GetHost(surl string) bool {
    oUrl, err := url.Parse(surl)
    if err != nil { return false }
    c.Host = oUrl.Host
    return true
}

func (c *Crawl) Crawler(surl string) {
    c.NewResources = c.NewResources[0:0]
    if strings.HasSuffix(surl,".html") || strings.HasSuffix(surl,".htm") { //TODO: js?
        if !c.IsCrawled(surl) && c.GetHost(surl) {
            c.Crawled = append(c.Crawled, surl)
            c.BaseUrl = surl
            c.Scan(surl)
            c.SubUrls()
        }
    }
}