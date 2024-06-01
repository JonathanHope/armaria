package messaging

import (
	"fmt"
	"io"
)

// The top level handler.
// Invokes the right handler for a given message kind.

// Dispatch branches the handler invocation based on message kind.
func Dispatch(reader io.Reader, writer io.Writer) error {
	in, err := ReceiveMessage(reader)
	if err != nil {
		return err
	}

	switch in.Kind {

	case MessageKindAddBook:
		if err := handleKind(writer, in, addBookHandler); err != nil {
			return err
		}

	case MessageKindAddFolder:
		if err := handleKind(writer, in, addFolderHandler); err != nil {
			return err
		}

	case MessageKindAddTags:
		if err := handleKind(writer, in, addTagsHandler); err != nil {
			return err
		}

	case MessageKindListBooks:
		if err := handleKind(writer, in, listBooksHandler); err != nil {
			return err
		}

	case MessageKindListTags:
		if err := handleKind(writer, in, listTagsHandler); err != nil {
			return err
		}

	case MessageKindRemoveBook:
		if err := handleKind(writer, in, removeBookHandler); err != nil {
			return err
		}

	case MessageKindRemoveFolder:
		if err := handleKind(writer, in, removeFolderHandler); err != nil {
			return err
		}

	case MessageKindRemoveTags:
		if err := handleKind(writer, in, removeTagsHandler); err != nil {
			return err
		}

	case MessageKindUpdateBook:
		if err := handleKind(writer, in, updateBookHandler); err != nil {
			return err
		}

	case MessageKindUpdateFolder:
		if err := handleKind(writer, in, updateFolderHandler); err != nil {
			return err
		}

	default:
		return fmt.Errorf("Unknown message kind: %s", in.Kind)
	}

	return nil
}

// handleKind invokes the right handler for a given message kind.
func handleKind(writer io.Writer, in NativeMessage, handler handleFn) error {
	out, err := handler(in)
	if err != nil {
		out, err = PayloadToMessage(MessageKindError, ErrorPayload{
			Error: fmt.Sprintf("%s", err),
		})
		if err != nil {
			return err
		}
	}

	err = out.SendMessage(writer)
	if err != nil {
		return err
	}

	return nil
}
