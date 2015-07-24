export GOPATH=$(shell pwd)

all:
	go get code.google.com/p/go.net/html
	go build dirscan.go crawl.go requests.go printer.go wordlist.go bruter.go config.go
	strip dirscan

linux32:
	GOOS=linux GOARCH=386 go build dirscan.go crawl.go requests.go printer.go wordlist.go bruter.go config.go

linux64:
	GOOS=linux GOARCH=amd64 go build dirscan.go crawl.go requests.go printer.go wordlist.go bruter.go config.go

win32:
	GOOS=windows GOARCH=386 go build dirscan.go crawl.go requests.go printer.go wordlist.go bruter.go config.go

win64:
 	GOOS=windows GOARCH=amd64 go build dirscan.go crawl.go requests.go printer.go wordlist.go bruter.go config.go

clean:
	rm -f dirscan *.exe
	
uninstall:
	rm -f /usr/bin/dirscan

install:
	cp dirscan /usr/bin/
