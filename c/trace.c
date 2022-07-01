#include "vmlinux.h"

#include <bpf/bpf_core_read.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

#define TASK_COMM_LEN 16

char __license[] SEC("license") = "Dual MIT/GPL";

struct event {
  u32 pid;
  s32 fd;
  u8 comm[TASK_COMM_LEN];
};
// Force emitting struct event into the ELF.
const struct event *unused_event __attribute__((unused));

// Change syscall to trace and check if `invalid bpf_context access` error occurs.
//
// SEC("fentry/__x64_sys_close")
SEC("fentry/__x64_sys_recvfrom")
int BPF_PROG(fentry_syscall, struct pt_regs *regs) {
  struct event t;

  bpf_get_current_comm(t.comm, TASK_COMM_LEN);

  u64 id = bpf_get_current_pid_tgid();
  t.pid = id >> 32;

  // Comment out this line to remove `invalid bpf_context access` error
  t.fd = PT_REGS_PARM1_CORE(regs);

  bpf_printk("comm: %s, pid: %d, fd: %d", t.comm, t.pid, t.fd);

  return 0;
}
