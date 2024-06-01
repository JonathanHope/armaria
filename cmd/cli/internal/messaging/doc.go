// messaging contains the logic for Armaria to communicate with JSON.
// This is the message format used to communicate with browser extensions.
// They are encoded by first writing the size of the message as a unit32.
// After that the message is encoded to JSON and written as binary.
// All of this should be done over stdout and stdin.
// This is also the format the JSON formatter uses.
// The messages require calls to unmarshal the JSON.
// First to get the kind of the message; second to unmarshal the payload once the type is known.
package messaging
