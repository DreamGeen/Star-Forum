// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: publish.proto

package publishPb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
	_ "star/proto/feed/feedPb"
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

// Api Endpoints for PublishService service

func NewPublishServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for PublishService service

type PublishService interface {
	PreUploadVideos(ctx context.Context, in *PreUploadVideosRequest, opts ...client.CallOption) (*PreUploadVideosResponse, error)
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...client.CallOption) (*CreatePostResponse, error)
	CountPost(ctx context.Context, in *CountPostRequest, opts ...client.CallOption) (*CountPostResponse, error)
	ListPost(ctx context.Context, in *ListPostRequest, opts ...client.CallOption) (*ListPostResponse, error)
}

type publishService struct {
	c    client.Client
	name string
}

func NewPublishService(name string, c client.Client) PublishService {
	return &publishService{
		c:    c,
		name: name,
	}
}

func (c *publishService) PreUploadVideos(ctx context.Context, in *PreUploadVideosRequest, opts ...client.CallOption) (*PreUploadVideosResponse, error) {
	req := c.c.NewRequest(c.name, "PublishService.PreUploadVideos", in)
	out := new(PreUploadVideosResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publishService) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...client.CallOption) (*CreatePostResponse, error) {
	req := c.c.NewRequest(c.name, "PublishService.CreatePost", in)
	out := new(CreatePostResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publishService) CountPost(ctx context.Context, in *CountPostRequest, opts ...client.CallOption) (*CountPostResponse, error) {
	req := c.c.NewRequest(c.name, "PublishService.CountPost", in)
	out := new(CountPostResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *publishService) ListPost(ctx context.Context, in *ListPostRequest, opts ...client.CallOption) (*ListPostResponse, error) {
	req := c.c.NewRequest(c.name, "PublishService.ListPost", in)
	out := new(ListPostResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PublishService service

type PublishServiceHandler interface {
	PreUploadVideos(context.Context, *PreUploadVideosRequest, *PreUploadVideosResponse) error
	CreatePost(context.Context, *CreatePostRequest, *CreatePostResponse) error
	CountPost(context.Context, *CountPostRequest, *CountPostResponse) error
	ListPost(context.Context, *ListPostRequest, *ListPostResponse) error
}

func RegisterPublishServiceHandler(s server.Server, hdlr PublishServiceHandler, opts ...server.HandlerOption) error {
	type publishService interface {
		PreUploadVideos(ctx context.Context, in *PreUploadVideosRequest, out *PreUploadVideosResponse) error
		CreatePost(ctx context.Context, in *CreatePostRequest, out *CreatePostResponse) error
		CountPost(ctx context.Context, in *CountPostRequest, out *CountPostResponse) error
		ListPost(ctx context.Context, in *ListPostRequest, out *ListPostResponse) error
	}
	type PublishService struct {
		publishService
	}
	h := &publishServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&PublishService{h}, opts...))
}

type publishServiceHandler struct {
	PublishServiceHandler
}

func (h *publishServiceHandler) PreUploadVideos(ctx context.Context, in *PreUploadVideosRequest, out *PreUploadVideosResponse) error {
	return h.PublishServiceHandler.PreUploadVideos(ctx, in, out)
}

func (h *publishServiceHandler) CreatePost(ctx context.Context, in *CreatePostRequest, out *CreatePostResponse) error {
	return h.PublishServiceHandler.CreatePost(ctx, in, out)
}

func (h *publishServiceHandler) CountPost(ctx context.Context, in *CountPostRequest, out *CountPostResponse) error {
	return h.PublishServiceHandler.CountPost(ctx, in, out)
}

func (h *publishServiceHandler) ListPost(ctx context.Context, in *ListPostRequest, out *ListPostResponse) error {
	return h.PublishServiceHandler.ListPost(ctx, in, out)
}
