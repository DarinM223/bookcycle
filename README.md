# Bookcycle

[![Build Status](https://travis-ci.org/DarinM223/bookcycle.svg)](https://travis-ci.org/DarinM223/bookcycle)

![Main Image](http://i.imgur.com/lgGEjsT.png)

Development setup
=================

First you have to install Go. You can do so by using Homebrew or downloading the distribution off of the Golang home page. Be sure to set the GOPATH environment variable to the root path of Go (most of the time it is in ~/go). Also the version of Go must support vendored dependencies (1.5+).

Then run:
```
go get github.com/DarinM223/bookcycle
```
Which should put the project inside $GOPATH/src/github.com/DarinM223/bookcycle. That is the folder where all of the development is done. 

Before pushing code to master, make sure that the code is formatted with gofmt (a formatting tool that comes with go). Plugins for many text editors like GoSublime for Sublime Text or vim-go for vim automatically do this whenever you save a Go file.

Building
=================
To get the latest version of the dependencies, first run 
```
go get
```

Then to build the project, first type
```
go install
```
to cache all non-main dependencies.

Then type
```
go build
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

Screenshots
===========

![New Book](http://i.imgur.com/jOqxuKI.png)

![Recent Books](http://i.imgur.com/ykuqznz.png)

![Book detail](http://i.imgur.com/XLeD4lu.png)

![Messaging](http://i.imgur.com/BIdsRdX.png)

