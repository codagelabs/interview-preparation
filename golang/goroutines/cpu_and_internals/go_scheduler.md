
# Go sheduler and how it works internally

Go sheduler is a part of the Go runtime that manages the execution of goroutines. It is responsible for scheduling goroutines to run on available threads and ensuring that they are executed efficiently.


## Key Components

The Go scheduler has three main components that work together to manage goroutine execution:

1. **Goroutines (G)**
   goroutines are Lightweight independently executing function (similar to a thread). managed by go runtime. It has its own stack, very small stack initially up 2kb, Can grow/shrink stack as needed


2. **OS Threads (M)**
The ‘M’ stands for machine. This Thread is still managed by the OS and the OS is still responsible for placing the Thread on a Core for execution. This means when I run a Go program on my machine, as pre system configuration no of threads available to execute my work, each individually attached to a P.
   

3. **Processors (P)**
   Logical processors that manage goroutine queues. They act as a bridge between M and G. Each P has its own local queue of goroutines. The number equals GOMAXPROCS setting. P is required for M to execute G.

The relationship between these components is often referred to as the "GPM model":
- Multiple Gs (goroutines) are scheduled onto
- P (processors) which are in turn scheduled onto
- M (OS threads)


## Goroutine states 
Just like Threads, Goroutines have the same three high-level states. A Goroutine can be in one of three states: Waiting, Runnable or Running/Executing.

1. **Running/Executing**
   - The goroutine is currently being executed by the Go scheduler
   - The goroutine is running on an M and using a P   

2. **Waiting**
   - The goroutine is waiting for something to complete (e.g. waiting for a channel operation, I/O, or a mutex)
   - The goroutine is not using any M or P

3. **Runnable**
   - The goroutine which are ready to run but not currently executing they are waiting to be executed by the Go scheduler
   - The goroutine is in a queue (local or global)
   - As more goroutines are created, the waiting time for each goroutine to get executed may increase due to queue length


## Context Switching

Context switching is the process of saving the current state of a running goroutine and restoring the state of a waiting goroutine. This is done to allow the scheduler to switch between goroutines on the same thread.
These are the points in the code where context switching happens:

1. **System Calls**
   - When a goroutine makes a system call (e.g. read, write, socket, etc.)
   - The goroutine is blocked and the scheduler switches to another goroutine 

2. **Channel Operations**
   - When a goroutine is waiting for a channel operation (e.g. receive, send)
   - The goroutine is blocked and the scheduler switches to another goroutine

3. **Garbage Collection**
   -GC has ist own set of goroutines to manage, when the garbage collector is running some goroutines may need to be stopped and started again, this is done by the scheduler.

4. **Go Statements**
   - When a goroutine is created using the `go` keyword, the goroutine is added to the local queue of the creating P  

5. **Blocking Operations**
   - When a goroutine is waiting for a blocking operation (e.g. I/O, mutex, etc.)
   - The goroutine is blocked and the scheduler switches to another goroutine    



## Cooperating Scheduler 

The Go scheduler is a cooperative scheduler that runs in user space as part of your application's runtime. Unlike the OS scheduler which is preemptive and runs in kernel space, the Go scheduler makes scheduling decisions at specific points in the code like:

- When making system/network calls
- During channel operations 
- When creating new goroutines
- During garbage collection
- During blocking operations

While it's technically cooperative, the Go scheduler behaves much like a preemptive scheduler in practice. This is because:

1. Scheduling decisions are handled by the runtime, not developers
2. The timing of switches is unpredictable
3. Developers don't need to explicitly yield control

This clever design gives Go programs the benefits of cooperative scheduling (efficiency, safety) while maintaining the feel and flexibility of preemptive scheduling.
## Key Terminology

1. **GOMAXPROCS**
   - Controls the maximum number of operating system threads that can execute Go code simultaneously
   - Usually set to match the number of CPU cores
   - Can be modified using runtime.GOMAXPROCS()
   - Default value is the number of CPU cores available

   ```
   runtime.GOMAXPROCS(runtime.NumCPU())
   // runtime.NumCPU() returns the number of CPU cores available
   // GOMAXPROCS is a runtime variable that controls the maximum number of operating system threads that can execute Go code simultaneously.
   runtime.GOMAXPROCS(1) // set to 1 thread
   runtime.GOMAXPROCS(2) // set to 2 thread
   ```

2. **Work Stealing**
   - Scheduling algorithm used by Go scheduler
   - When a P's local queue is empty, it tries to steal work from other P's queues
   - Helps balance load across all processors
   - Improves CPU utilization

3. **Run Queue**
   - Each P maintains a local run queue of goroutines
   - There's also a global run queue for the entire program
   - Local queues are checked first before the global queue
   - New goroutines are typically added to the local queue of the creating P

4. **Context Switch**
   - Process of storing and restoring execution state when switching between goroutines
   - Much lighter weight than OS thread context switches
   - Handled entirely by Go runtime
   - Can happen on channel operations, system calls, or when explicitly yielding

5. **Blocking Operations**
   - Operations that might cause a goroutine to wait (I/O, channel operations, system calls)
   - When a goroutine blocks, M can detach from P and execute other goroutines
   - Go runtime handles this efficiently to maintain concurrency
   - New M might be created if needed to keep P busy

6. **[Network Poller](network_poller_guide.md)**
   - Network poller is a component in Go's runtime that handles asynchronous I/O operations
   - Prevents blocking OS threads when performing network operations
   - Uses efficient system calls like epoll (Linux), kqueue (BSD/macOS), or IOCP (Windows)
   - Key responsibilities:
   * Monitors multiple network connections simultaneously
      


## How it works

The Go scheduler uses a work-stealing algorithm to distribute goroutines across multiple threads (known as OS threads or M). Here's how it works:

1. **M:N Scheduling Model**
   - Go implements an M:N scheduler where M goroutines are multiplexed onto N OS threads
   - The number of OS threads is typically tied to the number of CPU cores (GOMAXPROCS)

2. **Key Components**
   - G (Goroutine): The basic unit of execution
   - M (Machine): OS thread that can execute goroutines
   - P (Processor): A logical processor that manages a queue of goroutines

3. **Scheduling Process**
   - Each P has a local queue of runnable goroutines
   - When a P's queue is empty, it tries to steal work from other P's queues
   - If a goroutine makes a blocking call (like I/O), the M can switch to another goroutine
   - The scheduler ensures fair distribution of CPU time among goroutines

4. **Context Switching**
   - The scheduler performs context switches between goroutines very efficiently
   - Context switches are much cheaper than OS thread context switches
   - This allows Go to manage thousands of goroutines with minimal overhead

5. **Load Balancing**
   - Work stealing ensures even distribution of work across all available cores
   - If one P runs out of work, it can steal goroutines from other P's queues
   - This helps maintain CPU utilization and system performance


## Understanding Go's Runtime Scheduler

When running Go programs on my local machine with 8 CPU cores, here's how the runtime scheduler operates:

1. **Processors (P's)**
   - `runtime.NumCPU()` returns 8 on my machine
   - My Go programs automatically get 8 P's (logical processors)
   - Each P represents a scheduling context for running goroutines

2. **OS Threads (M's)**
   - Each P gets assigned one OS Thread (M for "machine")
   - The operating system manages these threads
   - The OS handles scheduling these threads onto CPU cores
   - With 8 P's, my program has 8 threads for executing work
   - Each thread (M) is individually attached to its own P

3. **Goroutines (G's)**
   - Every Go program starts with one initial Goroutine
   - A Goroutine is Go's version of a Coroutine (hence "G" instead of "C")
   - Think of Goroutines as application-level threads
   - Like OS threads switching on/off CPU cores
   - Goroutines switch on/off M's instead

4. **Run Queues**
   - The scheduler uses two types of run queues:
     * Global Run Queue (GRQ): Holds unassigned goroutines
     * Local Run Queue (LRQ): Each P has its own LRQ
   - Each P's LRQ manages goroutines for that P
   - Goroutines take turns being switched on/off the M assigned to their P
   - GRQ holds goroutines waiting to be assigned to a P's LRQ
   - (Process of moving goroutines from GRQ to LRQ will be covered later)

  


