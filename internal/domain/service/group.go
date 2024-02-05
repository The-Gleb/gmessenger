package service

import "github.com/The-Gleb/gmessenger/protos/gen/go/group"

type groupService struct {
	groupClient group.GroupClient
}

func NewGroupService(gc group.GroupClient) *groupService {
	return &groupService{
		groupClient: gc,
	}
}

func (gs *groupService) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {

}

func (gs *groupService) CheckMember(ctx context.Context, in *CheckMemberRequest, opts ...grpc.CallOption) (*CheckMembersResponse, error) {

}

func (gs *groupService) GetMembers(ctx context.Context, in *GetMembersRequest, opts ...grpc.CallOption) (*GetMembersResponse, error) {

}

func (gs *groupService) AddMessage(ctx context.Context, in *AddMessageRequest, opts ...grpc.CallOption) (*AddMessageResponse, error) {

}

func (gs *groupService) GetMessages(ctx context.Context, in *GetMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error) {

}

func (gs *groupService) UpdateMessageStatus(ctx context.Context, in *UpdateMessageStatusRequest, opts ...grpc.CallOption) (*UpdateMessageStatusResponse, error) {

}

func (gs *groupService) (ctx context.Context, in *GetGroupsRequest, opts ...grpc.CallOption) (*GetGroupsResponse, error) {

}

