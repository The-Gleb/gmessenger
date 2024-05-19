package serverapi

import (
	"context"
	"time"

	"github.com/The-Gleb/gmessenger/app/pkg/proto/go/group"
	"github.com/The-Gleb/gmessenger/group_service/gateway/domain/entity"
	"github.com/The-Gleb/gmessenger/group_service/gateway/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageService interface {
	AddMessage(ctx context.Context, message entity.Message) (entity.Message, error)
	GetMessages(ctx context.Context, groupID int64, limit, offset int) ([]entity.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID int64, status string) (entity.Message, error)
	GetLastMessage(ctx context.Context, groupID int64) (entity.Message, error)
}

type GroupService interface {
	Create(ctx context.Context, group entity.GroupCreate) (entity.Group, error)
	IsMember(ctx context.Context, userLogin string, groupID int64) (bool, error)
	GetMembers(ctx context.Context, groupID int64) ([]string, error)
	GetGroups(ctx context.Context, userLogin string, limit, offset int) ([]entity.Group, error)
	Exists(ctx context.Context, groupID int64) (bool, error)
}

type serverAPI struct {
	group.UnimplementedGroupServer
	messageService MessageService
	groupService   GroupService
}

func NewServerAPI(ms MessageService, gs GroupService) *serverAPI {
	return &serverAPI{
		messageService: ms,
		groupService:   gs,
	}
}

func (s *serverAPI) AddMessage(ctx context.Context, addMessageRequest *group.AddMessageRequest) (*group.AddMessageResponse, error) {
	if addMessageRequest.GetGroupId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}
	if addMessageRequest.GetSenderLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "sender login is required")
	}
	if addMessageRequest.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "message text is required")
	}

	exists, err := s.groupService.Exists(ctx, addMessageRequest.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn`t check if group exists")
	}

	if !exists {
		return nil, status.Error(codes.NotFound, "group not found")
	}

	newMessage, err := s.messageService.AddMessage(ctx, entity.Message{
		Sender:    addMessageRequest.GetSenderLogin(),
		Text:      addMessageRequest.GetText(),
		GroupID:   addMessageRequest.GetGroupId(),
		Status:    entity.SENT,
		Timestamp: time.Now(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn`t add new message")
	}

	return &group.AddMessageResponse{
		Message: &group.MessageView{
			Id:          newMessage.ID,
			GroupId:     newMessage.GroupID,
			SenderLogin: newMessage.Sender,
			Text:        newMessage.Text,
			Status:      group.MessageStatus(group.MessageStatus_value[newMessage.Status]),
			Timestamp:   timestamppb.New(newMessage.Timestamp),
		},
	}, nil

}
func (s *serverAPI) CheckMember(ctx context.Context, checkMemberRequest *group.CheckMemberRequest) (*group.CheckMembersResponse, error) {

	if checkMemberRequest.GetGroupId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}
	if checkMemberRequest.GetUserLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "user login is required")
	}

	isMember, err := s.groupService.IsMember(ctx, checkMemberRequest.GetUserLogin(), checkMemberRequest.GetGroupId())
	if err != nil {
		if errors.Code(err) == errors.ErrGroupNotFound {
			return nil, status.Error(codes.InvalidArgument, "group doesn`t exists")
		}
		return nil, status.Error(codes.Internal, "couldn`t check if user is a memeber of the group")
	}

	return &group.CheckMembersResponse{
		IsMember: isMember,
	}, nil

}
func (s *serverAPI) Create(ctx context.Context, createRequst *group.CreateRequest) (*group.CreateResponse, error) {

	if createRequst.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "group name is required")
	}
	if createRequst.GetMemberLogins() == nil || len(createRequst.GetMemberLogins()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "members` logins are required")
	}

	newGroup, err := s.groupService.Create(ctx, entity.GroupCreate{
		Name:          createRequst.GetName(),
		MembersLogins: createRequst.GetMemberLogins(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "error adding group to db")
	}

	return &group.CreateResponse{
		Id: newGroup.ID,
	}, nil

}
func (s *serverAPI) GetGroups(ctx context.Context, getGroupsRequest *group.GetGroupsRequest) (*group.GetGroupsResponse, error) {

	if getGroupsRequest.GetUserLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "group name is required")
	}
	if getGroupsRequest.GetLimit() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message limit is required")
	}

	groups, err := s.groupService.GetGroups(
		ctx,
		getGroupsRequest.GetUserLogin(),
		int(getGroupsRequest.GetLimit()),
		int(getGroupsRequest.GetOffset()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "error getting user's groups")
	}

	groupViews := make([]*group.GroupView, len(groups))

	for i, g := range groups {

		groupViews[i] = &group.GroupView{
			Id:     g.ID,
			Name:   g.Name,
			Unread: 0, // TODO: figure out how to store unread messages for each group member
		}

		lastMessage, err := s.messageService.GetLastMessage(ctx, g.ID)
		if err != nil {
			if errors.Code(err) == errors.ErrNoDataFound {
				groupViews[i].LastMessage = nil
				continue
			}
			return nil, status.Error(codes.Internal, "error getting user's groups")
		}

		groupViews[i].LastMessage = &group.MessageView{
			Id:          lastMessage.ID,
			GroupId:     lastMessage.GroupID,
			SenderLogin: lastMessage.Sender,
			Text:        lastMessage.Text,
			Status:      group.MessageStatus(group.MessageStatus_value[lastMessage.Status]),
			Timestamp:   timestamppb.New(lastMessage.Timestamp),
		}

	}

	return &group.GetGroupsResponse{
		Groups: groupViews,
	}, nil

}
func (s *serverAPI) GetMembers(ctx context.Context, getMembersRequest *group.GetMembersRequest) (*group.GetMembersResponse, error) {

	if getMembersRequest.GetGroupId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}

	members, err := s.groupService.GetMembers(ctx, getMembersRequest.GetGroupId())
	if err != nil {
		if errors.Code(err) == errors.ErrGroupNotFound {
			return nil, status.Error(codes.NotFound, "group not found")
		}
		return nil, status.Error(codes.Internal, "group id is required")
	}

	return &group.GetMembersResponse{
		MemberLogins: members,
	}, nil

}
func (s *serverAPI) GetMessages(ctx context.Context, getMessagesRequest *group.GetMessagesRequest) (*group.GetMessagesResponse, error) {

	if getMessagesRequest.GetGroupId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "group id is required")
	}

	exists, err := s.groupService.Exists(ctx, getMessagesRequest.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn`t check if group exists")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "group not found")
	}

	messages, err := s.messageService.GetMessages(
		ctx,
		getMessagesRequest.GetGroupId(),
		int(getMessagesRequest.GetLimit()),
		int(getMessagesRequest.GetOffset()),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "couldn`t get group messages")
	}

	grpcMessages := make([]*group.MessageView, len(messages))

	for i, msg := range messages {

		grpcMessages[i] = &group.MessageView{
			Id:          msg.ID,
			SenderLogin: msg.Sender,
			Status:      group.MessageStatus(group.MessageStatus_value[msg.Status]),
			Text:        msg.Text,
			GroupId:     msg.GroupID,
			Timestamp:   timestamppb.New(msg.Timestamp),
		}
	}

	return &group.GetMessagesResponse{
		Messages: grpcMessages,
	}, nil

}
func (s *serverAPI) UpdateMessageStatus(ctx context.Context, updateMessageStatusRequest *group.UpdateMessageStatusRequest) (*group.UpdateMessageStatusResponse, error) {

	if updateMessageStatusRequest.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "message id is required")
	}
	if updateMessageStatusRequest.GetStatus().String() == "" {
		return nil, status.Error(codes.InvalidArgument, "message status is required")
	}

	msg, err := s.messageService.UpdateMessageStatus(
		ctx,
		updateMessageStatusRequest.GetId(),
		updateMessageStatusRequest.Status.String(),
	)
	if err != nil {
		if errors.Code(err) == errors.ErrNoDataFound {
			return nil, status.Error(codes.InvalidArgument, "message not found")
		}
		return nil, status.Error(codes.Internal, "couldn`t update group messages")
	}

	return &group.UpdateMessageStatusResponse{
		Id:     msg.ID,
		Status: group.MessageStatus(group.MessageStatus_value[msg.Status]),
	}, nil

}
