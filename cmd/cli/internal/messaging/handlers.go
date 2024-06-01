package messaging

import (
	"github.com/jonathanhope/armaria/pkg"
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

	options := armaria.DefaultAddBookOptions()
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

	book, err := armaria.AddBook(payload.URL, options)
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

	options := armaria.DefaultAddFolderOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	if payload.ParentID.Valid {
		options.WithParentID(payload.ParentID.String)
	}

	book, err := armaria.AddFolder(payload.Name, options)
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

	options := armaria.DefaultAddTagsOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}
	book, err := armaria.AddTags(payload.ID, payload.Tags, options)
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

	options := armaria.DefaultListBooksOptions()
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

	books, err := armaria.ListBooks(options)
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

	options := armaria.DefaultListTagsOptions()
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

	tags, err := armaria.ListTags(options)
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

	options := armaria.DefaultRemoveBookOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}

	err = armaria.RemoveBook(payload.ID, options)
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

	options := armaria.DefaultRemoveFolderOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}

	err = armaria.RemoveFolder(payload.ID, options)
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

	options := armaria.DefaultRemoveTagsOptions()
	if payload.DB.Valid {
		options.WithDB(payload.DB.String)
	}

	book, err := armaria.RemoveTags(payload.ID, payload.Tags, options)
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

	options := armaria.DefaultUpdateBookOptions()
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
	if payload.NextBook.Valid {
		options.WithOrderBefore(payload.NextBook.String)
	}
	if payload.PreviousBook.Valid {
		options.WithOrderAfter(payload.PreviousBook.String)
	}

	book, err := armaria.UpdateBook(payload.ID, options)
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

	options := armaria.DefaultUpdateFolderOptions()
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
	if payload.NextBook.Valid {
		options.WithOrderBefore(payload.NextBook.String)
	}
	if payload.PreviousBook.Valid {
		options.WithOrderAfter(payload.PreviousBook.String)
	}

	book, err := armaria.UpdateFolder(payload.ID, options)
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
