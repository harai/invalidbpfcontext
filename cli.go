//go:build linux
// +build linux

package main

import "github.com/harai/invalidbpfcontext/trace"

func main() {
	trace.Run()
}
