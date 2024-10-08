// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: comment.proto

package commentPb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for CommentService service

func NewCommentServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for CommentService service

type CommentService interface {
	// 发布评论服务
	PostComment(ctx context.Context, in *PostCommentRequest, opts ...client.CallOption) (*PostCommentResponse, error)
	// 删除评论服务
	DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...client.CallOption) (*DeleteCommentResponse, error)
	// 获取评论服务
	GetComments(ctx context.Context, in *GetCommentsRequest, opts ...client.CallOption) (*GetCommentsResponse, error)
	CountComment(ctx context.Context, in *CountCommentRequest, opts ...client.CallOption) (*CountCommentResponse, error)
	QueryCommentExist(ctx context.Context, in *QueryCommentExistRequest, opts ...client.CallOption) (*QueryCommentExistResponse, error)
}

type commentService struct {
	c    client.Client
	name string
}

func NewCommentService(name string, c client.Client) CommentService {
	return &commentService{
		c:    c,
		name: name,
	}
}

func (c *commentService) PostComment(ctx context.Context, in *PostCommentRequest, opts ...client.CallOption) (*PostCommentResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.PostComment", in)
	out := new(PostCommentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentService) DeleteComment(ctx context.Context, in *DeleteCommentRequest, opts ...client.CallOption) (*DeleteCommentResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.DeleteComment", in)
	out := new(DeleteCommentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentService) GetComments(ctx context.Context, in *GetCommentsRequest, opts ...client.CallOption) (*GetCommentsResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.GetComments", in)
	out := new(GetCommentsResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentService) CountComment(ctx context.Context, in *CountCommentRequest, opts ...client.CallOption) (*CountCommentResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.CountComment", in)
	out := new(CountCommentResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentService) QueryCommentExist(ctx context.Context, in *QueryCommentExistRequest, opts ...client.CallOption) (*QueryCommentExistResponse, error) {
	req := c.c.NewRequest(c.name, "CommentService.QueryCommentExist", in)
	out := new(QueryCommentExistResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CommentService service

type CommentServiceHandler interface {
	// 发布评论服务
	PostComment(context.Context, *PostCommentRequest, *PostCommentResponse) error
	// 删除评论服务
	DeleteComment(context.Context, *DeleteCommentRequest, *DeleteCommentResponse) error
	// 获取评论服务
	GetComments(context.Context, *GetCommentsRequest, *GetCommentsResponse) error
	CountComment(context.Context, *CountCommentRequest, *CountCommentResponse) error
	QueryCommentExist(context.Context, *QueryCommentExistRequest, *QueryCommentExistResponse) error
}

func RegisterCommentServiceHandler(s server.Server, hdlr CommentServiceHandler, opts ...server.HandlerOption) error {
	type commentService interface {
		PostComment(ctx context.Context, in *PostCommentRequest, out *PostCommentResponse) error
		DeleteComment(ctx context.Context, in *DeleteCommentRequest, out *DeleteCommentResponse) error
		GetComments(ctx context.Context, in *GetCommentsRequest, out *GetCommentsResponse) error
		CountComment(ctx context.Context, in *CountCommentRequest, out *CountCommentResponse) error
		QueryCommentExist(ctx context.Context, in *QueryCommentExistRequest, out *QueryCommentExistResponse) error
	}
	type CommentService struct {
		commentService
	}
	h := &commentServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&CommentService{h}, opts...))
}

type commentServiceHandler struct {
	CommentServiceHandler
}

func (h *commentServiceHandler) PostComment(ctx context.Context, in *PostCommentRequest, out *PostCommentResponse) error {
	return h.CommentServiceHandler.PostComment(ctx, in, out)
}

func (h *commentServiceHandler) DeleteComment(ctx context.Context, in *DeleteCommentRequest, out *DeleteCommentResponse) error {
	return h.CommentServiceHandler.DeleteComment(ctx, in, out)
}

func (h *commentServiceHandler) GetComments(ctx context.Context, in *GetCommentsRequest, out *GetCommentsResponse) error {
	return h.CommentServiceHandler.GetComments(ctx, in, out)
}

func (h *commentServiceHandler) CountComment(ctx context.Context, in *CountCommentRequest, out *CountCommentResponse) error {
	return h.CommentServiceHandler.CountComment(ctx, in, out)
}

func (h *commentServiceHandler) QueryCommentExist(ctx context.Context, in *QueryCommentExistRequest, out *QueryCommentExistResponse) error {
	return h.CommentServiceHandler.QueryCommentExist(ctx, in, out)
}
