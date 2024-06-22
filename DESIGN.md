# MessageTracker Implementation Design

The type `tracker` in the package `network` implements the `network.MessageTracker` interface.

The first design decision to point out is that I slightly changed the API of the network package, specifically the
`NewMessageTracker()` contructor function. This was done to handle length values less than or equal to zero as errors.
Depending on the context, changing this API might not be feasible. In that case, the check can be performed in `Add()`
instead.

As far as other edge cases are concerned, `Add()` considers `nil` messages and messages with an empty ID as invalid.
The existing tests already provided 100% coverage, except for these edge cases, for which I added tests as well.

Internally messages are stored in a doubly-linked list and in a map, with the message ID used as key. This data
structure was chosen to ensure time and space complexity of O(1) for `Add()`, `Message()` and `Delete()`. The design
was guided by benchmarks, starting out with a simple implementation that used a slice for storing messages. Time and
space complexity of `Messages()` is O(n), which I consider acceptable since it is assumed to be called less frequently
and in less performance-critical code paths.
