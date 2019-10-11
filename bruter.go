package main

import (
	"net/url"
	"strings"
)

//import "fmt"

type BruterEndCallback func(res *[]string)

type Bruter struct {
	Url        string
	Chan       chan string
	Resources  []string
	Routines   int
	EndCB      BruterEndCallback
	Extensions []string
	Magic      string
}

func (b *Bruter) InjectWords() {
	for _, w := range W.Words {
		b.Chan <- w
	}
	close(b.Chan)
}

func (b *Bruter) OnEnd(end BruterEndCallback) {
	b.EndCB = end
}

func (b *Bruter) Brute(surl string) {
	b.Chan = make(chan string, 6)
	b.Extensions = CFG.Extensions
	b.Routines = CFG.Goroutines
	b.Url = surl
	b.Magic = "check1337"

	//fmt.Println("bruteforcing ",surl)

	go b.InjectWords()

	for i := 0; i < b.Routines; i++ {
		go b.Worker(surl, i)
	}
}

func (b *Bruter) IsRepeated(surl string) bool {
	for _, u := range b.Resources {
		if u == surl {
			return true
		}
	}
	return false
}

func (b *Bruter) Check(surl string, ext string) {
	_, err := url.Parse(surl)
	if err != nil {
		return
	}

	//TODO: refactor this
	if ext == "" {
		html, code, _ := R.Get(surl + "/")
		if code != 404 {
			if !b.IsRepeated(surl + ext) {
				P.Show("b", code, len(html), surl+ext)
				b.Resources = append(b.Resources, surl+ext)
			}
		}
	}

	html, code, _ := R.Get(surl + ext)
	words := len(strings.Split(html, " "))

	if SzList.Push(words) {
		if !b.IsRepeated(surl + ext) {
			P.Show("b", code, len(html), surl+ext)
			b.Resources = append(b.Resources, surl+ext)
		}
	}
}

func (b *Bruter) Worker(surl string, r int) {

	for w := range b.Chan {
		b.Check(surl+w, "")

		for _, x := range b.Extensions {
			b.Check(surl+w, "."+x)
		}
	}

	b.Routines--
	if b.Routines == 0 && len(b.Resources) > 0 {
		b.EndCB(&b.Resources)
	}

}
