package websocket

import (
	"star/constant/str"
	"sync"
)

type HubManager struct {
	member sync.Map
}

func NewHubManager() *HubManager {
	return &HubManager{}
}

func (h *HubManager) CreateHub(communityId int64) error {
	_, ok := h.member.Load(communityId)
	if !ok {
		newHub := NewHub()
		h.member.Store(communityId, newHub)
		go newHub.Run()
		return nil
	}
	return str.ErrHubExists
}

func (h *HubManager) GetHub(communityId int64) (*Hub, error) {
	value, ok := h.member.Load(communityId)
	if !ok {
		return nil, str.ErrHubNotExists
	}
	return value.(*Hub), nil
}

func (h *HubManager) DeleteHub(communityId int64) {
	h.member.Delete(communityId)
}

func (h *HubManager) Run() {
	h.member.Range(func(key, value interface{}) bool {
		go value.(*Hub).Run()
		return true
	})
}
