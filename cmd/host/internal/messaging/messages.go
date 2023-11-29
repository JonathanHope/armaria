package messaging

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

// This is the message format used to communicate with browser extensions.
// They are encoded by first writing the size of the message as a unit32.
// After that the message is encoded to JSON and written as binary.
// All of this should be done over stdout and stdin.

// MessageKind denotes the kind of message that was sent to or received from a browser extension.
type MessageKind string

const (
	MessageKindError        MessageKind = "error"         // message contains an error that occurred
	MessageKindBooks        MessageKind = "books"         // message contains zero or more books
	MessageKindBook         MessageKind = "book"          // message contains a single book
	MessageKindVoid         MessageKind = "void"          // message contains nothing
	MessageKindTags         MessageKind = "tags"          // message contains zero or more tags
	MessageKindAddBook      MessageKind = "add-book"      // message is a request to add a bookmark
	MessageKindAddFolder    MessageKind = "add-folder"    // message is a request to add a bookmark
	MessageKindAddTags      MessageKind = "add-tags"      // message is a request to add tags
	MessageKindListBooks    MessageKind = "list-books"    // message is a request to list books
	MessageKindListTags     MessageKind = "list-tags"     // message is a request to list tags
	MessageKindRemoveBook   MessageKind = "remove-book"   // message is a request to remove a bookmark
	MessageKindRemoveFolder MessageKind = "remove-folder" // message is a request to remove a folder
	MessageKindRemoveTags   MessageKind = "remove-tags"   // message is a request to remove tags from a bookmark
	MessageKindUpdateBook   MessageKind = "update-book"   // message is a request to update a bookmark
	MessageKindUpdateFolder MessageKind = "update-folder" // message is a request to update a folder
)

// NativeMessage is a message sent to or received from a browser extension.
type NativeMessage struct {
	Kind    MessageKind `json:"kind"`    // denotes what kind of message this is
	Payload string      `json:"payload"` // a JSON payload that is different depending on the MessageKind
}

// SendMessage sends a message to a browser extension.
func (msg NativeMessage) SendMessage(writer io.Writer) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.LittleEndian, uint32(len(messageBytes)))
	if err != nil {
		return err
	}

	_, err = writer.Write(messageBytes)
	if err != nil {
		return err
	}

	return nil
}

// ReceiveMessage receives a message from a browser extension.
func ReceiveMessage(reader io.Reader) (NativeMessage, error) {
	var messageLength uint32
	err := binary.Read(reader, binary.LittleEndian, &messageLength)
	if err != nil {
		return NativeMessage{}, err
	}

	messageBytes := make([]byte, messageLength)
	_, err = reader.Read(messageBytes)
	if err != nil {
		return NativeMessage{}, err
	}

	var message NativeMessage
	err = json.Unmarshal(messageBytes, &message)
	if err != nil {
		return NativeMessage{}, err
	}

	return message, nil
}
