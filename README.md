# piano-server

A simple service to perform MIDI input/output, optimized for piano, written in [Ell](https://github.com/boynton/ell).

This app is an extension of [midi-ell](https://github.com/boynton/midi-ell), itself an extension to Ell that
adds native MIDI I/O. This app is also an example of a program written in Ell, but made into an executable via
the standard Go installation method.

## Running the server

	$ go get github.com/boynton/piano-server
	$ piano-server

Then, you can connect to port 8888, and send framed messages to make midi play. See piano-client.ell for examples
of that, it is a file that can be loaded into a stock Ell runtime (no extensions needed).

## Protocol

A TCP socket is listened on (port selected at launch). For new connections accepted, communication is
framed. A single varint defined the number of bytes in the frame, followed by those bytes. The varint
is zig-zag encoded as a signed int (as in Avro, protobuf, and Go's encoding/binary use). The payload
of the frame in this example is a JSON array where the first element is the operation, and other elements are
the parameters of the operation:

	["note" key vel dur]
	["sustain" amount]

The key, vel, and amount numbers are between 0 and 127, and the dur is specified in floating point seconds.

The note-offs are scheduled the duration time into the future, and are guaranteed to happen even if the
client closed the connection before scheduled time.

