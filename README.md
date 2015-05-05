# Bookcycle

### Project for CS130

Neil Bedi – 404018872

Han Lee - 804011454

Darin Minamoto – 704140102

Rebecca Pan – 603929588

Belinda Yang- 604021996

Development setup
=================

First you have to install Go. You can do so by using Homebrew or downloading the distribution off of the Golang home page. Be sure to set the GOPATH environment variable to the root path of Go (most of the time it is in ~/go). 

Then run:
```
go get github.com/DarinM223/cs130-test
```
Which should put the project inside $GOPATH/src/github.com/DarinM223/cs130-test. That is the folder where all of the development is done. 

Before pushing code to master, make sure that the code is formatted with gofmt (a formatting tool that comes with go). Plugins for many text editors like GoSublime for Sublime Text or vim-go for vim automatically do this whenever you save a Go file.

Running
=======

To run the server, type
```
go run *.go
```
into the terminal at the project path. Navigate to [http://localhost:8080](http://localhost:8080) and it should display the home page.