## IO

https://www.educative.io/answers/how-to-read-and-write-with-golang-bufio

In Golang, bufio is a package used for buffered IO. Buffering IO is a technique used to temporarily accumulate the results for an IO operation before transmitting it forward. This technique can increase the speed of a program by reducing the number of system calls, which are typically slow operations.

User buffered I/O, shortened to buffering or buffered I/O, refers to the technique of temporarily storing the results of an I/O operation in user-space before transmitting it to the kernel (in the case of writes) or before providing it to your process (in the case of reads). By so buffering the data, you can minimize the number of system calls and can block-align I/O operations, which may improve the performance of your application.

For example, consider a process that writes one character at a time to a file. This is obviously inefficient: Each write operation corresponds to a write() system call, which means a trip into the kernel, a memory copy (of a single byte!), and a return to user-space, only to repeat the whole ordeal. Worse, filesystems and storage media work in terms of blocks; operations are fastest when aligned to integer multiples of those blocks. Misaligned operations, particularly very small ones, incur additional overhead.

User buffered I/O avoids this inefficiency by buffering the writes in a data buffer in user-space until a certain threshold is reached—ideally, the underlying filesystem's block size or an integer multiple thereof. To use our previous example, we'd simply copy each character into the buffer and call write() only when the block size is reached.

A similar process happens with reads. Imagine a process that reads one line of a file into memory at a time. That might be logical given how the program works, but not optimal given how the filesystem and underlying storage media work. Thus, user buffering could read a large, block-aligned chunk of the file into a buffer, and then dole out small pieces of it as requested by the process.

producer --> buffer --> io.Writer

buffer has space for 4 characters
producer     buffer          destination (io.Writer)
a            ----->   a
b            ----->   ab
c            ----->   abc
d            ----->   abcd
e            ----->   e      ------>   abcd
f            ----->   ef               abcd
g            ----->   efg              abcd
h            ----->   efgh             abcd
i            ----->   i      ------>   abcdefgh