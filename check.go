package main

import (
	"fmt"
	"net"
)

func checkHostAvailable(host string) bool {
	conn, err := net.Dial("tcp", host)

	if nil != err {
		fmt.Println("Error", err)
		return false
	}

	conn.Close()
	return true
}
