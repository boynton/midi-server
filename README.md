# piano-server

A simple service to perform MIDI input/output, optimized for piano.

This is an example of using Ell embedded in a an app. This app is equivalent to Ell, but
defines primitives to talk to MIDI.

## Protocol

A TCP socket is listened on (port selected at launch). For new connections accepted, communication is
framed. A single varint defined the number of bytes in the frame, followed by those bytes. The payload
of the frame is currently JSON array where the first element is the operation, and other elements are
the parameters of the operation:

   ["note" key vel dur]
   ["sustain" amount]

The note-offs are scheduled the duration time into the future, and are guaranteed to happen even if the
client closed the connection before scheduled time.


      
