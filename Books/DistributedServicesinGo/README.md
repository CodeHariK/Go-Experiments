# Distributed Services in go

## Log : A log is an append-only sequence of records

Logs—which are sometimes also called write-ahead logs, transaction logs, or commit logs—are at the heart of storage engines, message queues, version control, and replication and consensus algorithms. 

Folks who develop storage engines of filesystems and databases use logs to
improve the data integrity of their systems. The ext filesystems, for example,
log changes to a journal instead of directly changing the disk’s data file. Once
the filesystem has safely written the changes to the journal, it then applies
those changes to the data files. Logging to the journal is simple and fast, so
there’s little chance of losing data. Even if your computer crashed before ext
had finished updating the disk files, then on the next boot, the filesystem
would process the data in the journal to complete its updates. Database
developers, like PostgreSQL, use the same technique to make their systems
durable: they record changes to a log, called a write-ahead log (WAL), and
later process the WAL to apply the changes to their database’s data files.

Database developers use the WAL for replication, too. Instead of writing the
logs to a disk, they write the logs over the network to its replicas. The replicas
apply the changes to their own data copies, and eventually they all end up
at the same state. Raft, a consensus algorithm, uses the same idea to get
distributed services to agree on a cluster-wide state. Each node in a Raft
cluster runs a state machine with a log as its input. The leader of the Raft
cluster appends changes to its followers’ logs. Since the state machines use
the logs as input and because the logs have the same records in the same
order, all the services end up with the same state.

Distributed version control systems like Git work similarly; your commit log history is a literal commit log. As you can see, a complete log not only holds the latest state, but all states that have existed.

A log is like a table that always orders the records by time and indexes each record by its
offset and time created. Concrete implementations of logs have to deal with us not having disks with infinite space, which means we can’t append to the same file forever. So we
split the log into a list of segments. When the log grows too big, we free up
disk space by deleting old segments whose data we’ve already processed or
archived. This cleaning up of old segments can run in a background process
while our service can still produce to the active (newest) segment and consume
from other segments with no, or at least fewer, conflicts where goroutines
access the same data.
There’s always one special segment among the list of segments, and that’s
the active segment. We call it the active segment because it’s the only segment
we actively write to. When we’ve filled the active segment, we create a new
segment and make it the active segment.

Each segment comprises a store file and an index file. The segment’s store
file is where we store the record data; we continually append records to this
file. The segment’s index file is where we index each record in the store file.
The index file speeds up reads because it maps record offsets to their position
in the store file. Reading a record given its offset is a two-step process: first
you get the entry from the index file for the record, which tells you the position
of the record in the store file, and then you read the record at that position
in the store file. Since the index file requires only two small fields—the offset
and stored position of the record—the index file is much smaller than the
store file that stores all your record data. Index files are small enough that we can memory-map them and make operations on the file as fast as operating on in-memory data.

Record—the data stored in our log.
Store—the file we store records in.
Index—the file we store index entries in.
Segment—the abstraction that ties a store and an index together.
Log—the abstraction that ties all the segments together.