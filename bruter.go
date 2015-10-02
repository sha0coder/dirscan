package main

import "net/url"

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

	_, code_not_found, _ := R.Get(surl + b.Magic + ext)
	//fmt.Printf("not_found:%d\n", code_not_found)
	html, code, _ := R.Get(surl + ext)
	if code != code_not_found && code > 0 {

		if !b.IsRepeated(surl + ext) { //TODO: meter dentro de este if toda esta funcion
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
