# `invalid bpf_context access` error on some distributions

`invalid bpf_context access` error occurs on some distributions when trying to read parameters passed to the syscall traced with fentry/fexit:

```c
t.fd = PT_REGS_PARM1_CORE(regs);
```

This seems to depend on how `struct pt_regs *regs` is defined in `/sys/kernel/btf/vmlinux` file.

While this example traces `recvfrom` syscall, you can change it whatever you like to compare the results.

## Prerequisites

* Go 1.18
* Latest Clang and LLVM
* Latest Libbpf
* `bpftool`

## Build

```
script/generate-vmlinux-header
go generate trace/trace.go
go build -o ./output .
```

## Run

```
sudo ./output
```

## Comparison

### Amazon Linux 2 (with `kernel-5.15` package)

```
$ uname -a
Linux ip-10-1-1-66.ap-northeast-1.compute.internal 5.15.43-20.123.amzn2.x86_64 #1 SMP Fri May 27 00:28:44 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux
```

#### BTF `__x64_sys_recvmsg` entry

```
$ bpftool btf dump file /sys/kernel/btf/vmlinux format raw
...
[15439] FWD 'pt_regs' fwd_kind=struct
[15440] CONST '(anon)' type_id=15439
[15441] PTR '(anon)' type_id=15440
[15442] FUNC_PROTO '(anon)' ret_type_id=34 vlen=1
        '__unused' type_id=15441
...
[15694] FUNC '__x64_sys_recvmsg' type_id=15442 linkage=static
...
```

#### Result

Error occurred.

```
$ sudo ./output
2022/07/01 12:16:54 loading objects: field FentrySyscall: program fentry_syscall: load program: permission denied: invalid bpf_context access off=0 size=8 (5 line(s) omitted)
```

### Ubuntu 22.04 LTS

```
$ uname -a
Linux ip-10-1-1-30 5.15.0-1014-aws #18-Ubuntu SMP Wed Jun 15 20:04:04 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux
```

#### BTF `__x64_sys_recvmsg` entry

```
$ bpftool btf dump file /sys/kernel/btf/vmlinux format raw
...
[17252] FWD 'pt_regs' fwd_kind=struct
[17253] CONST '(anon)' type_id=17252
[17254] PTR '(anon)' type_id=17253
[17255] FUNC_PROTO '(anon)' ret_type_id=35 vlen=1
        '__unused' type_id=17254
...
[17519] FUNC '__x64_sys_recvmsg' type_id=17255 linkage=static
...
```

#### Result

Error occurred.

```
$ sudo ./output
2022/07/01 03:25:43 loading objects: field FentrySyscall: program fentry_syscall: load program: permission denied: invalid bpf_context access off=0 size=8 (5 line(s) omitted)
```

### Debian 11

```
$ uname -a
Linux ip-10-1-1-101 5.10.0-14-cloud-amd64 #1 SMP Debian 5.10.113-1 (2022-04-29) x86_64 GNU/Linux
```

#### BTF `__x64_sys_recvmsg` entry

```
$ bpftool btf dump file /sys/kernel/btf/vmlinux format raw
...
[13362] FWD 'pt_regs' fwd_kind=struct
[13363] CONST '(anon)' type_id=13362
[13364] PTR '(anon)' type_id=13363
[13365] FUNC_PROTO '(anon)' ret_type_id=36 vlen=1
        '__unused' type_id=13364
...
[13608] FUNC '__x64_sys_recvmsg' type_id=13365 linkage=static
...
```

#### Result

Error occurred.

```
$ sudo ./output
2022/07/01 03:33:01 loading objects: field FentrySyscall: program fentry_syscall: load program: permission denied: invalid bpf_context access off=0 size=8 (6 line(s) omitted)
```

### Ubuntu 20.04 LTS (with newer kernel)

```
$ uname -a
Linux xxx-XPS-13-9300 5.13.0-51-generic #58~20.04.1-Ubuntu SMP Tue Jun 14 11:29:12 UTC 2022 x86_64 x86_64 x86_64 GNU/Linux
```

#### BTF `__x64_sys_recvmsg` entry

```
$ bpftool btf dump file /sys/kernel/btf/vmlinux format raw
[1] INT 'long unsigned int' size=8 bits_offset=0 nr_bits=64 encoding=(none)
...
[226] STRUCT 'pt_regs' size=168 vlen=21
        'r15' type_id=1 bits_offset=0
        'r14' type_id=1 bits_offset=64
        'r13' type_id=1 bits_offset=128
        'r12' type_id=1 bits_offset=192
        'bp' type_id=1 bits_offset=256
        'bx' type_id=1 bits_offset=320
        'r11' type_id=1 bits_offset=384
        'r10' type_id=1 bits_offset=448
        'r9' type_id=1 bits_offset=512
        'r8' type_id=1 bits_offset=576
        'ax' type_id=1 bits_offset=640
        'cx' type_id=1 bits_offset=704
        'dx' type_id=1 bits_offset=768
        'si' type_id=1 bits_offset=832
        'di' type_id=1 bits_offset=896
        'orig_ax' type_id=1 bits_offset=960
        'ip' type_id=1 bits_offset=1024
        'cs' type_id=1 bits_offset=1088
        'flags' type_id=1 bits_offset=1152
        'sp' type_id=1 bits_offset=1216
        'ss' type_id=1 bits_offset=1280
...
[5183] CONST '(anon)' type_id=226
...
[5189] PTR '(anon)' type_id=5183
...
[5321] FUNC_PROTO '(anon)' ret_type_id=42 vlen=1
        '__unused' type_id=5189
...
[17648] FUNC '__x64_sys_recvmsg' type_id=5321 linkage=static
...
```

#### Result

No errors occur.

```
$ sudo ./output
2022/07/01 12:53:32 Open /sys/kernel/debug/tracing/trace_pipe to read events.

$ sudo cat /sys/kernel/debug/tracing/trace_pipe
       IPCServer-6778    [003] d... 18280.241302: bpf_trace_printk: comm: IPCServer, pid: 6776, fd: 11
 ibus-engine-moz-2982    [001] d... 18280.241332: bpf_trace_printk: comm: ibus-engine-moz, pid: 2982, fd: 206
 ibus-engine-moz-2982    [001] d... 18280.241337: bpf_trace_printk: comm: ibus-engine-moz, pid: 2982, fd: 206
 APEX_CONTEXT_NE-2948    [000] d... 18280.505341: bpf_trace_printk: comm: APEX_CONTEXT_NE, pid: 2918, fd: 96
```
