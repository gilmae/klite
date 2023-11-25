# What is this?

I have an idea for an embedded kafka. Sqlite is an embedded relational database with a sql engine. Embedded meaning you add a dependency/lib and it can be used in code to access a relational database. Same thing, but for an event journal like Kafka.

klite. Dekaf. I can't decide what to call it

It's kind of a dumb idea, yeah, but it's an idea I want to play with because why not.

## The Plan
Intention is to be a (mostly) append-only stream of data, that can later be retrieved by key and in chunks. The query language might look something like:

`ADD $data TO $stream`

`GET $key FROM $stream`

`GET $key[, $key2[, $key3]] FROM $stream`

`GET $num AFTER $key FROM $stream`

`GET $num BEFORE $key FROM $stream`

### Currently we have:

1. A linked list of nodes that acts as the value store.
2. A b-tree index of keys. Each node value in the b-tree points to:
    1. A page in the linked list
    2. The offset within the page where the value starts
    3. A length. The data can span multiple nodes in the linked list, which are pages of 4096 bytes.
3. Functions to add new sets of data to the stream
4. Functions to retrieve data from the stream by key
5. Functions to retrieve n items from the stream starting with key x

### What we think we need
* To help support the BEFORE and AFTER commands, Next and Previous links in the value headers
* A repl (plus command parser)
* Support for multiple streams. Not sure how to store a hash of string to stream root page in the file. Another B Tree?

More long range things:

* Transactions
* Expiring items
