# Design Data Intensive Application

* https://github.com/ept/ddia-references

Security
Efficiency
Consistency
Reliability, Fault Tolerance, Resilient : Replication, Testing
Scalability : Load, Performance (Network, Disk, CPU, Memory)
Observability
Maintainibility

## Data in at center of many challenges

• Store data so that they, or another application, can find it again later (databases)
• Remember the result of an expensive operation, to speed up reads (caches)
• Allow users to search data by keyword or filter it in various ways (search indexes)
• Send a message to another process, to be handled asynchronously (stream pro‐
cessing)
• Periodically crunch a large amount of accumulated data (batch processing)

## Database

* Bitcask
* LSM Tree [Writes], SS Table, Memtable
* B-Tree [Reads]
