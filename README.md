# Bookcycle

[![Build Status](https://travis-ci.org/DarinM223/bookcycle)](https://travis-ci.org/DarinM223/bookcycle)

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
go get github.com/DarinM223/bookcycle
```
Which should put the project inside $GOPATH/src/github.com/DarinM223/bookcycle. That is the folder where all of the development is done. 

Before pushing code to master, make sure that the code is formatted with gofmt (a formatting tool that comes with go). Plugins for many text editors like GoSublime for Sublime Text or vim-go for vim automatically do this whenever you save a Go file.

Building normally
=================
To get the latest version of the dependencies, first run 
```
go get
```

Then to run the server, first type
```
go build
```
to build the application into the bookcycle executable. 

Building with godep
===================
This project also has godep versioned dependencies. If you want to use those instead of getting the latest version, after installing godep with 
```
go get github.com/tools/godep
```

Then run
```
godep go build
```
to build the application into the bookcycle executable.

Running
=======
Enter
```
./bookcycle
```
into the terminal at the project path to run the application. Navigate to [http://localhost:8080](http://localhost:8080) and it should display the home page.

You might run into an error named something like this when logging in:
```
gob: type not registered for interface *main.User
```
If this happens, you have to clear out your cookies by going into the chrome development tools and clicking Resource and the dropdown arrow under Cookies. There should be a localhost option. Right click that and click Clear to c lear the cookies. After that login should work

Documentation
=============
Documentation for all methods used for the backend is in https://godoc.org/github.com/DarinM223/bookcycle/server
