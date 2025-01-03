// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: like.proto

package likePb

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

// Api Endpoints for LikeService service

func NewLikeServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for LikeService service

type LikeService interface {
	LikeAction(ctx context.Context, in *LikeActionRequest, opts ...client.CallOption) (*LikeActionResponse, error)
	GetUserTotalLike(ctx context.Context, in *GetUserTotalLikeRequest, opts ...client.CallOption) (*GetUserTotalLikeResponse, error)
	LikeList(ctx context.Context, in *LikeListRequest, opts ...client.CallOption) (*LikeListResponse, error)
	GetLikeCount(ctx context.Context, in *GetLikeCountRequest, opts ...client.CallOption) (*GetLikeCountResponse, error)
	IsLike(ctx context.Context, in *IsLikeRequest, opts ...client.CallOption) (*IsLikeResponse, error)
	GetUserLikeCount(ctx context.Context, in *GetUserLikeCountRequest, opts ...client.CallOption) (*GetUserLikeCountResponse, error)
}

type likeService struct {
	c    client.Client
	name string
}

func NewLikeService(name string, c client.Client) LikeService {
	return &likeService{
		c:    c,
		name: name,
	}
}

func (c *likeService) LikeAction(ctx context.Context, in *LikeActionRequest, opts ...client.CallOption) (*LikeActionResponse, error) {
	req := c.c.NewRequest(c.name, "LikeService.LikeAction", in)
	out := new(LikeActionResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeService) GetUserTotalLike(ctx context.Context, in *GetUserTotalLikeRequest, opts ...client.CallOption) (*GetUserTotalLikeResponse, error) {
	req := c.c.NewRequest(c.name, "LikeService.GetUserTotalLike", in)
	out := new(GetUserTotalLikeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeService) LikeList(ctx context.Context, in *LikeListRequest, opts ...client.CallOption) (*LikeListResponse, error) {
	req := c.c.NewRequest(c.name, "LikeService.LikeList", in)
	out := new(LikeListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeService) GetLikeCount(ctx context.Context, in *GetLikeCountRequest, opts ...client.CallOption) (*GetLikeCountResponse, error) {
	req := c.c.NewRequest(c.name, "LikeService.GetLikeCount", in)
	out := new(GetLikeCountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeService) IsLike(ctx context.Context, in *IsLikeRequest, opts ...client.CallOption) (*IsLikeResponse, error) {
	req := c.c.NewRequest(c.name, "LikeService.IsLike", in)
	out := new(IsLikeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likeService) GetUserLikeCount(ctx context.Context, in *GetUserLikeCountRequest, opts ...client.CallOption) (*GetUserLikeCountResponse, error) {
	req := c.c.NewRequest(c.name, "LikeService.GetUserLikeCount", in)
	out := new(GetUserLikeCountResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LikeService service

type LikeServiceHandler interface {
	LikeAction(context.Context, *LikeActionRequest, *LikeActionResponse) error
	GetUserTotalLike(context.Context, *GetUserTotalLikeRequest, *GetUserTotalLikeResponse) error
	LikeList(context.Context, *LikeListRequest, *LikeListResponse) error
	GetLikeCount(context.Context, *GetLikeCountRequest, *GetLikeCountResponse) error
	IsLike(context.Context, *IsLikeRequest, *IsLikeResponse) error
	GetUserLikeCount(context.Context, *GetUserLikeCountRequest, *GetUserLikeCountResponse) error
}

func RegisterLikeServiceHandler(s server.Server, hdlr LikeServiceHandler, opts ...server.HandlerOption) error {
	type likeService interface {
		LikeAction(ctx context.Context, in *LikeActionRequest, out *LikeActionResponse) error
		GetUserTotalLike(ctx context.Context, in *GetUserTotalLikeRequest, out *GetUserTotalLikeResponse) error
		LikeList(ctx context.Context, in *LikeListRequest, out *LikeListResponse) error
		GetLikeCount(ctx context.Context, in *GetLikeCountRequest, out *GetLikeCountResponse) error
		IsLike(ctx context.Context, in *IsLikeRequest, out *IsLikeResponse) error
		GetUserLikeCount(ctx context.Context, in *GetUserLikeCountRequest, out *GetUserLikeCountResponse) error
	}
	type LikeService struct {
		likeService
	}
	h := &likeServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&LikeService{h}, opts...))
}

type likeServiceHandler struct {
	LikeServiceHandler
}

func (h *likeServiceHandler) LikeAction(ctx context.Context, in *LikeActionRequest, out *LikeActionResponse) error {
	return h.LikeServiceHandler.LikeAction(ctx, in, out)
}

func (h *likeServiceHandler) GetUserTotalLike(ctx context.Context, in *GetUserTotalLikeRequest, out *GetUserTotalLikeResponse) error {
	return h.LikeServiceHandler.GetUserTotalLike(ctx, in, out)
}

func (h *likeServiceHandler) LikeList(ctx context.Context, in *LikeListRequest, out *LikeListResponse) error {
	return h.LikeServiceHandler.LikeList(ctx, in, out)
}

func (h *likeServiceHandler) GetLikeCount(ctx context.Context, in *GetLikeCountRequest, out *GetLikeCountResponse) error {
	return h.LikeServiceHandler.GetLikeCount(ctx, in, out)
}

func (h *likeServiceHandler) IsLike(ctx context.Context, in *IsLikeRequest, out *IsLikeResponse) error {
	return h.LikeServiceHandler.IsLike(ctx, in, out)
}

func (h *likeServiceHandler) GetUserLikeCount(ctx context.Context, in *GetUserLikeCountRequest, out *GetUserLikeCountResponse) error {
	return h.LikeServiceHandler.GetUserLikeCount(ctx, in, out)
}
