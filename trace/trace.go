//go:build linux
// +build linux

package trace

import (
	"errors"
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
		if err2 := errors.Unwrap(err); err2 != nil {
			if err3 := errors.Unwrap(err2); err3 != nil {
				if err4 := errors.Unwrap(err3); err4 != nil {
					log.Printf("loading objects: %+v", err4)
				}
				log.Printf("loading objects: %+v", err3)
			}
			log.Printf("loading objects: %+v", err2)
		}
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
