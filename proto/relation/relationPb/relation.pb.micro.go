// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: relation.proto

package relationPb

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
	_ "star/proto/user/userPb"
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

// Api Endpoints for RelationService service

func NewRelationServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for RelationService service

type RelationService interface {
	GetFollowerList(ctx context.Context, in *GetFollowRequest, opts ...client.CallOption) (*GetFollowResponse, error)
	GetFansList(ctx context.Context, in *GetFansRequest, opts ...client.CallOption) (*GetFansResponse, error)
	CountFollower(ctx context.Context, in *CountFollowRequest, opts ...client.CallOption) (*CountFollowResponse, error)
	CountFans(ctx context.Context, in *CountFansRequest, opts ...client.CallOption) (*CountFansResponse, error)
	Follow(ctx context.Context, in *FollowRequest, opts ...client.CallOption) (*FollowResponse, error)
	UnFollow(ctx context.Context, in *UnFollowRequest, opts ...client.CallOption) (*UnFollowResponse, error)
	IsFollow(ctx context.Context, in *IsFollowRequest, opts ...client.CallOption) (*IsFollowResponse, error)
}

type relationService struct {
	c    client.Client
	name string
}

func NewRelationService(name string, c client.Client) RelationService {
	return &relationService{
		c:    c,
		name: name,
	}
}

func (c *relationService) GetFollowerList(ctx context.Context, in *GetFollowRequest, opts ...client.CallOption) (*GetFollowResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.GetFollowerList", in)
	out := new(GetFollowResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) GetFansList(ctx context.Context, in *GetFansRequest, opts ...client.CallOption) (*GetFansResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.GetFansList", in)
	out := new(GetFansResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) CountFollower(ctx context.Context, in *CountFollowRequest, opts ...client.CallOption) (*CountFollowResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.CountFollower", in)
	out := new(CountFollowResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) CountFans(ctx context.Context, in *CountFansRequest, opts ...client.CallOption) (*CountFansResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.CountFans", in)
	out := new(CountFansResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) Follow(ctx context.Context, in *FollowRequest, opts ...client.CallOption) (*FollowResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.Follow", in)
	out := new(FollowResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) UnFollow(ctx context.Context, in *UnFollowRequest, opts ...client.CallOption) (*UnFollowResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.UnFollow", in)
	out := new(UnFollowResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relationService) IsFollow(ctx context.Context, in *IsFollowRequest, opts ...client.CallOption) (*IsFollowResponse, error) {
	req := c.c.NewRequest(c.name, "RelationService.IsFollow", in)
	out := new(IsFollowResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RelationService service

type RelationServiceHandler interface {
	GetFollowerList(context.Context, *GetFollowRequest, *GetFollowResponse) error
	GetFansList(context.Context, *GetFansRequest, *GetFansResponse) error
	CountFollower(context.Context, *CountFollowRequest, *CountFollowResponse) error
	CountFans(context.Context, *CountFansRequest, *CountFansResponse) error
	Follow(context.Context, *FollowRequest, *FollowResponse) error
	UnFollow(context.Context, *UnFollowRequest, *UnFollowResponse) error
	IsFollow(context.Context, *IsFollowRequest, *IsFollowResponse) error
}

func RegisterRelationServiceHandler(s server.Server, hdlr RelationServiceHandler, opts ...server.HandlerOption) error {
	type relationService interface {
		GetFollowerList(ctx context.Context, in *GetFollowRequest, out *GetFollowResponse) error
		GetFansList(ctx context.Context, in *GetFansRequest, out *GetFansResponse) error
		CountFollower(ctx context.Context, in *CountFollowRequest, out *CountFollowResponse) error
		CountFans(ctx context.Context, in *CountFansRequest, out *CountFansResponse) error
		Follow(ctx context.Context, in *FollowRequest, out *FollowResponse) error
		UnFollow(ctx context.Context, in *UnFollowRequest, out *UnFollowResponse) error
		IsFollow(ctx context.Context, in *IsFollowRequest, out *IsFollowResponse) error
	}
	type RelationService struct {
		relationService
	}
	h := &relationServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&RelationService{h}, opts...))
}

type relationServiceHandler struct {
	RelationServiceHandler
}

func (h *relationServiceHandler) GetFollowerList(ctx context.Context, in *GetFollowRequest, out *GetFollowResponse) error {
	return h.RelationServiceHandler.GetFollowerList(ctx, in, out)
}

func (h *relationServiceHandler) GetFansList(ctx context.Context, in *GetFansRequest, out *GetFansResponse) error {
	return h.RelationServiceHandler.GetFansList(ctx, in, out)
}

func (h *relationServiceHandler) CountFollower(ctx context.Context, in *CountFollowRequest, out *CountFollowResponse) error {
	return h.RelationServiceHandler.CountFollower(ctx, in, out)
}

func (h *relationServiceHandler) CountFans(ctx context.Context, in *CountFansRequest, out *CountFansResponse) error {
	return h.RelationServiceHandler.CountFans(ctx, in, out)
}

func (h *relationServiceHandler) Follow(ctx context.Context, in *FollowRequest, out *FollowResponse) error {
	return h.RelationServiceHandler.Follow(ctx, in, out)
}

func (h *relationServiceHandler) UnFollow(ctx context.Context, in *UnFollowRequest, out *UnFollowResponse) error {
	return h.RelationServiceHandler.UnFollow(ctx, in, out)
}

func (h *relationServiceHandler) IsFollow(ctx context.Context, in *IsFollowRequest, out *IsFollowResponse) error {
	return h.RelationServiceHandler.IsFollow(ctx, in, out)
}
