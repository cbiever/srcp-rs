# Welcome to SRCP-RS

SRCP-RS provides a REST interface for SRCP servers. It is a sister project to [Gember](https://github.com/cbiever/gember) which is a web application to control a model railroad.

## Installation

It is assumed a SRCP server, e.g. [SRCPD](http://srcpd.sourceforge.net/srcpd/index.html), is running.

 - Install Go: https://golang.org/doc/install
 - Clone the srcp-rs repositiory: git clone https://github.com/cbiever/srcp-rs
 - build with: go build

If **srcp-rs** is started without arguments it connects to port 4303 on localhost. Otherwise the first argument is used to connect to the SRCP server. 
