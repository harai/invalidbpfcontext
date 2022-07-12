//go:build linux
// +build linux

package trace

import (
	"log"
	"time"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror -D__TARGET_ARCH_x86" -type event bpf ../c/trace.c -- -I../c/headers

func Run() {
	// Allow the current process to lock memory for eBPF resources.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal(err)
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %+v", err)
	}
	defer objs.Close()

	tracing, err := link.AttachTracing(link.TracingOptions{Program: objs.FentrySyscall})
	if err != nil {
		log.Fatalf("attaching tracing: %+v", err)
	}
	defer tracing.Close()

	log.Println("Open /sys/kernel/debug/tracing/trace_pipe to read events.")
	time.Sleep(1 * time.Hour)
}
