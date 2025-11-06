package main

import "testing"

func TestProtocol(t *testing.T) {
	msg := "*3\r\n$3\r\nset\r\n$6\r\nleader\r\n$7\r\nCharlie\r\n"
	parseCommand(msg)
}
