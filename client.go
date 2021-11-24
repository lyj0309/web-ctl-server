package main

import (
	"net"
)

var ConnPool []con

type con struct {
	Ip    string
	Name  string
	State string
	Con   *net.TCPConn
}
