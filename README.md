# What is this?
An implementation of a btree

## Why is this?
I have an idea for an embedded kafka. Sqlite is an embedded relational database with a sql engine. Embedded meaning you add a dependency/lib and it can be used in code to access a relational database. Same thing, but for an event journal like Kafka.

It's kind of a dumb idea, yeah, but it's an idea I want to play with because why not.

So it starts with a btree. And because this btree is a means to an end, it's not all generic. The values this btree are holding are very specific. They are two uint32s. Exactly what those uint32s _mean_ is a little unclear to me at the moment. I am pretty sure the first is a page number in the file. The second might be a length of bytes. It might be a length of pages. The idea being that the messages being stored will be stored as blobs, and the values in the tree just point off to where the blob corresponding to a key can be found.

Which I am pretty sure is how sqlite does it.