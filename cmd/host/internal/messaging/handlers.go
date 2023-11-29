package messaging

import (
	"github.com/jonathanhope/armaria"
	"github.com/jonathanhope/armaria/pkg/api"
	"github.com/samber/lo"
)

// handleFn is a function that can handle a particular kind of message.
type handleFn func(in NativeMessage) (NativeMessage, error)

// addBookHandler handles an add-book message.
func addBookHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[AddBookPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultAddBookOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	if payload.Name.Valid {
		options.WithName(payload.Name.String)
	}
	if payload.Description.Valid {
		options.WithDescription(payload.Description.String)
	}
	if payload.ParentID.Valid {
		options.WithParentID(payload.ParentID.String)
	}
	if payload.Tags != nil {
		options.WithTags(payload.Tags)
	}

	book, err := armariaapi.AddBook(payload.URL, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBook, BookPayload{
		Book: bookMapper(book),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// addFolderHandler handles an add-folder message.
func addFolderHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[AddFolderPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultAddFolderOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	if payload.ParentID.Valid {
		options.WithParentID(payload.ParentID.String)
	}

	book, err := armariaapi.AddFolder(payload.Name, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBook, BookPayload{
		Book: bookMapper(book),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// addTagsHandler handles an add-tags message.
func addTagsHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[AddTagsPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultAddTagsOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	book, err := armariaapi.AddTags(payload.ID, payload.Tags, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBook, BookPayload{
		Book: bookMapper(book),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// listBooksHandler handles a list-books message.
func listBooksHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[ListBooksPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultListBooksOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	options.WithFolders(payload.IncludeFolders)
	options.WithBooks(payload.IncludeBookmarks)
	if payload.ParentID.Valid {
		options.WithParentID(payload.ParentID.String)
	}
	if payload.WithoutParentID {
		options.WithoutParentID()
	}
	if payload.After.Valid {
		options.WithAfter(payload.After.String)
	}
	if payload.Query.Valid {
		options.WithQuery(payload.Query.String)
	}
	if payload.Tags != nil {
		options.WithTags(payload.Tags)
	}
	if payload.Order != "" {
		options.WithOrder(armaria.Order(payload.Order))
	}
	if payload.Direction != "" {
		options.WithDirection(armaria.Direction(payload.Direction))
	}
	if payload.First.Valid {
		options.WithFirst(payload.First.Int64)
	}

	books, err := armariaapi.ListBooks(options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBooks, BooksPayload{
		Books: lo.Map(books, func(book armaria.Book, _ int) BookDTO {
			return bookMapper(book)
		}),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// listTagsHandler handles a list-tags message.
func listTagsHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[ListTagsPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultListTagsOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	if payload.After.Valid {
		options.WithAfter(payload.After.String)
	}
	if payload.Query.Valid {
		options.WithQuery(payload.Query.String)
	}
	if payload.Direction != "" {
		options.WithDirection(armaria.Direction(payload.Direction))
	}
	if payload.First.Valid {
		options.WithFirst(payload.First.Int64)
	}

	tags, err := armariaapi.ListTags(options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindTags, TagsPayload{
		Tags: tags,
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// removeBookHandler handles a remove-book message.
func removeBookHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[RemoveBookPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultRemoveBookOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}

	err = armariaapi.RemoveBook(payload.ID, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindVoid, VoidPayload{})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// removeFolderHandler handles a remove-book message.
func removeFolderHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[RemoveFolderPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultRemoveFolderOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}

	err = armariaapi.RemoveFolder(payload.ID, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindVoid, VoidPayload{})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// removeTagsHandler handles a remove-tags message.
func removeTagsHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[RemoveTagsPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultRemoveTagsOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}

	book, err := armariaapi.RemoveTags(payload.ID, payload.Tags, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBook, BookPayload{
		Book: bookMapper(book),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// updateBookHandler handles an update-book message.
func updateBookHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[UpdateBookPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultUpdateBookOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	if payload.RemoveDescription {
		options.WithoutDescription()
	}
	if payload.Description.Valid {
		options.WithDescription(payload.Description.String)
	}
	if payload.RemoveParentID {
		options.WithoutParentID()
	}
	if payload.ParentID.Valid {
		options.WithParentID(payload.ParentID.String)
	}
	if payload.Name.Valid {
		options.WithName(payload.Name.String)
	}
	if payload.URL.Valid {
		options.WithURL(payload.URL.String)
	}

	book, err := armariaapi.UpdateBook(payload.ID, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBook, BookPayload{
		Book: bookMapper(book),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}

// updateFolderHandler handles an update-folder message.
func updateFolderHandler(in NativeMessage) (NativeMessage, error) {
	payload, err := GetPayload[UpdateFolderPayload](in)
	if err != nil {
		return NativeMessage{}, err
	}

	options := armariaapi.DefaultUpdateFolderOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	if payload.RemoveParentID {
		options.WithoutParentID()
	}
	if payload.ParentID.Valid {
		options.WithParentID(payload.ParentID.String)
	}
	if payload.Name.Valid {
		options.WithName(payload.Name.String)
	}

	book, err := armariaapi.UpdateFolder(payload.ID, options)
	if err != nil {
		return NativeMessage{}, err
	}

	out, err := PayloadToMessage(MessageKindBook, BookPayload{
		Book: bookMapper(book),
	})
	if err != nil {
		return NativeMessage{}, err
	}

	return out, nil
}
