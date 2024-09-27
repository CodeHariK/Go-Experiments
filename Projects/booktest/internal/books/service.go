// Code generated by sqlc-connect (https://github.com/walterwanderley/sqlc-connect). DO NOT EDIT.

package books

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"

	pb "booktest/api/books/v1"
	"booktest/api/books/v1/v1connect"
	"booktest/internal/validation"
)

type Service struct {
	v1connect.UnimplementedBooksServiceHandler
	querier *Queries
}

func (s *Service) BooksByTags(ctx context.Context, req *connect.Request[pb.BooksByTagsRequest]) (*connect.Response[pb.BooksByTagsResponse], error) {
	dollar_1 := req.Msg.GetDollar_1()

	result, err := s.querier.BooksByTags(ctx, dollar_1)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "BooksByTags")
		return nil, err
	}
	res := new(pb.BooksByTagsResponse)
	for _, r := range result {
		res.List = append(res.List, toBooksByTagsRow(r))
	}
	return connect.NewResponse(res), nil
}

func (s *Service) BooksByTitleYear(ctx context.Context, req *connect.Request[pb.BooksByTitleYearRequest]) (*connect.Response[pb.BooksByTitleYearResponse], error) {
	var arg BooksByTitleYearParams
	arg.Title = req.Msg.GetTitle()
	arg.Year = req.Msg.GetYear()

	result, err := s.querier.BooksByTitleYear(ctx, arg)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "BooksByTitleYear")
		return nil, err
	}
	res := new(pb.BooksByTitleYearResponse)
	for _, r := range result {
		res.List = append(res.List, toBook(r))
	}
	return connect.NewResponse(res), nil
}

func (s *Service) CreateAuthor(ctx context.Context, req *connect.Request[pb.CreateAuthorRequest]) (*connect.Response[pb.CreateAuthorResponse], error) {
	name := req.Msg.GetName()

	result, err := s.querier.CreateAuthor(ctx, name)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "CreateAuthor")
		return nil, err
	}
	return connect.NewResponse(&pb.CreateAuthorResponse{Author: toAuthor(result)}), nil
}

func (s *Service) CreateBook(ctx context.Context, req *connect.Request[pb.CreateBookRequest]) (*connect.Response[pb.CreateBookResponse], error) {
	var arg CreateBookParams
	arg.AuthorID = req.Msg.GetAuthorId()
	arg.Isbn = req.Msg.GetIsbn()
	arg.BookType = BookType(req.Msg.GetBookType())
	arg.Title = req.Msg.GetTitle()
	arg.Year = req.Msg.GetYear()
	if v := req.Msg.GetAvailable(); v != nil {
		if err := v.CheckValid(); err != nil {
			err = fmt.Errorf("invalid Available: %s%w", err.Error(), validation.ErrUserInput)
			return nil, err
		}
		arg.Available = v.AsTime()
	} else {
		err := fmt.Errorf("field Available is required%w", validation.ErrUserInput)
		return nil, err
	}
	arg.Tags = req.Msg.GetTags()

	result, err := s.querier.CreateBook(ctx, arg)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "CreateBook")
		return nil, err
	}
	return connect.NewResponse(&pb.CreateBookResponse{Book: toBook(result)}), nil
}

func (s *Service) DeleteBook(ctx context.Context, req *connect.Request[pb.DeleteBookRequest]) (*connect.Response[pb.DeleteBookResponse], error) {
	bookID := req.Msg.GetBookId()

	err := s.querier.DeleteBook(ctx, bookID)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "DeleteBook")
		return nil, err
	}
	return connect.NewResponse(&pb.DeleteBookResponse{}), nil
}

func (s *Service) GetAuthor(ctx context.Context, req *connect.Request[pb.GetAuthorRequest]) (*connect.Response[pb.GetAuthorResponse], error) {
	authorID := req.Msg.GetAuthorId()

	result, err := s.querier.GetAuthor(ctx, authorID)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "GetAuthor")
		return nil, err
	}
	return connect.NewResponse(&pb.GetAuthorResponse{Author: toAuthor(result)}), nil
}

func (s *Service) GetBook(ctx context.Context, req *connect.Request[pb.GetBookRequest]) (*connect.Response[pb.GetBookResponse], error) {
	bookID := req.Msg.GetBookId()

	result, err := s.querier.GetBook(ctx, bookID)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "GetBook")
		return nil, err
	}
	return connect.NewResponse(&pb.GetBookResponse{Book: toBook(result)}), nil
}

func (s *Service) UpdateBook(ctx context.Context, req *connect.Request[pb.UpdateBookRequest]) (*connect.Response[pb.UpdateBookResponse], error) {
	var arg UpdateBookParams
	arg.Title = req.Msg.GetTitle()
	arg.Tags = req.Msg.GetTags()
	arg.BookType = BookType(req.Msg.GetBookType())
	arg.BookID = req.Msg.GetBookId()

	err := s.querier.UpdateBook(ctx, arg)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "UpdateBook")
		return nil, err
	}
	return connect.NewResponse(&pb.UpdateBookResponse{}), nil
}

func (s *Service) UpdateBookISBN(ctx context.Context, req *connect.Request[pb.UpdateBookISBNRequest]) (*connect.Response[pb.UpdateBookISBNResponse], error) {
	var arg UpdateBookISBNParams
	arg.Title = req.Msg.GetTitle()
	arg.Tags = req.Msg.GetTags()
	arg.BookID = req.Msg.GetBookId()
	arg.Isbn = req.Msg.GetIsbn()

	err := s.querier.UpdateBookISBN(ctx, arg)
	if err != nil {
		slog.Error("sql call failed", "error", err, "method", "UpdateBookISBN")
		return nil, err
	}
	return connect.NewResponse(&pb.UpdateBookISBNResponse{}), nil
}