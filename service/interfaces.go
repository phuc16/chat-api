package service

import (
	"app/dto"
	"app/entity"
	"app/repository"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -source $GOFILE -destination ../mocks/$GOPACKAGE/mock_$GOFILE -package mocks

type IUserRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveUser(ctx context.Context, user *entity.User) (err error)
	FindUserByID(ctx context.Context, id string) (*entity.User, error)
	AppendFriendRequest(ctx context.Context, id string, friendRequest entity.FriendRequest) (*mongo.UpdateResult, error)
	RemoveFriendRequest(ctx context.Context, senderID, receiver string) (*mongo.UpdateResult, error)
	UpdateTypeConversation(ctx context.Context, senderID, chatID1, chatID2, convoType string) (*mongo.UpdateResult, error)
	AppendConversation(ctx context.Context, id string, conversation entity.Conversation) (*mongo.UpdateResult, error)
	AppendConversationToMultiple(ctx context.Context, ids []string, conversation entity.Conversation) (*mongo.UpdateResult, error)
	RemoveConversation(ctx context.Context, id, chatID string) (*mongo.UpdateResult, error)
	RemoveConversationFromMultiple(ctx context.Context, ids []string, chatID string) (*mongo.UpdateResult, error)
	SearchConversation(ctx context.Context, senderID, chatID1, chatID2 string) (*entity.User, error)
	SearchSingleConversation(ctx context.Context, senderID, chatID string) (*entity.User, error)
	UpdateChatActivity(ctx context.Context, chatID string, lastUpdateAt time.Time, deliveries, reads []entity.Delivery, newTopChatActivity []entity.ChatActivity) (*mongo.UpdateResult, error)
	UpdateAvatarInConversation(ctx context.Context, userID, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInConversation(ctx context.Context, userID, newName string) (*mongo.UpdateResult, error)
	UpdateAvatarInFriendRequest(ctx context.Context, userID, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInFriendRequest(ctx context.Context, userID, newName string) (*mongo.UpdateResult, error)
	UpdateChatNameInConversation(ctx context.Context, ids []string, chatID, chatName string) (*mongo.UpdateResult, error)
	UpdateAvatarInConversationMultiple(ctx context.Context, ids []string, chatID, newAvatar string) (*mongo.UpdateResult, error)
	DeleteUserByID(ctx context.Context, id string) error
	GetAllRecentSearchProfiles(ctx context.Context, userID string) ([]entity.Profile, error)
	UpdateRecentSearchProfiles(ctx context.Context, userID string, recentSearchProfiles []entity.Profile) error
}

type ITokenRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	CreateToken(ctx context.Context, token *entity.Token) error
	GetTokenByID(ctx context.Context, id string) (*entity.Token, error)
	DeleteToken(ctx context.Context, token *entity.Token) error
}

type IAccountRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveAccount(ctx context.Context, account *entity.Account) error
	FindAccountByID(ctx context.Context, id string) (*entity.Account, error)
	GetAllAccounts(ctx context.Context) ([]*entity.Account, error)
	SearchByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.Account, error)
	ChangePassword(ctx context.Context, phoneNumber, password string) (*mongo.UpdateResult, error)
	ChangeAvatar(ctx context.Context, phoneNumber string, profile entity.Profile) (*mongo.UpdateResult, error)
	SearchByUserID(ctx context.Context, userID string) (*entity.Account, error)
	DeleteAccountByID(ctx context.Context, id string) error
	UpdateAccount(ctx context.Context, account *entity.Account) error
}

type IChatRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveChat(ctx context.Context, chat *entity.Chat) error
	FindChatByID(ctx context.Context, id string) (*entity.Chat, error)
	DeleteChatByID(ctx context.Context, chatID string) error
	AppendChatActivityByIDChat(ctx context.Context, chatID string, chatActivity entity.ChatActivity) (*mongo.UpdateResult, error)
	AppendDelivery(ctx context.Context, chatID string, delivery entity.Delivery) (*mongo.UpdateResult, error)
	AppendRead(ctx context.Context, chatID string, delivery entity.Delivery) (*mongo.UpdateResult, error)
	SearchDeliveryByUserID(ctx context.Context, chatID string, userID string) (*entity.Chat, error)
	SearchReadByUserID(ctx context.Context, chatID string, userID string) (*entity.Chat, error)
	ChangeDelivery(ctx context.Context, chatID string, userID string, messageID string) (*mongo.UpdateResult, error)
	ChangeRead(ctx context.Context, chatID string, userID string, messageID string) (*mongo.UpdateResult, error)
	RemoveDelivery(ctx context.Context, chatID string, messageID string) (*mongo.UpdateResult, error)
	RemoveRead(ctx context.Context, chatID string, messageID string) (*mongo.UpdateResult, error)
	AppendHiddenMessage(ctx context.Context, chatID string, userID string, messageID string) (*mongo.UpdateResult, error)
	RecallMessage(ctx context.Context, chatID string, messageID string) (*mongo.UpdateResult, error)
	GetChatTop10(ctx context.Context, chatID string) (*entity.Chat, error)
	GetChatActivityFromNToM(ctx context.Context, chatID string, x int, y int) ([]entity.ChatActivity, error)
	UpdateAvatarInRead(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInRead(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error)
	UpdateAvatarInDelivery(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInDelivery(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error)
	SearchByKeyWord(ctx context.Context, chatID string, key string) ([]entity.ChatActivity, error)
	GetIndexOfMessageID(ctx context.Context, chatID string, messageID string) ([]repository.SearchIndexes, error)
}

type IGroupRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveGroup(ctx context.Context, group *entity.Group) error
	FindGroupByID(ctx context.Context, id string) (*entity.Group, error)
	AppendMember(ctx context.Context, id string, info entity.PersonInfo) (*mongo.UpdateResult, error)
	RemoveMember(ctx context.Context, id string, userID string) (*mongo.UpdateResult, error)
	AppendAdmin(ctx context.Context, id string, info entity.PersonInfo) (*mongo.UpdateResult, error)
	RemoveAdmin(ctx context.Context, id string, userID string) (*mongo.UpdateResult, error)
	UpdateNameChat(ctx context.Context, id string, chatName string) (*mongo.UpdateResult, error)
	UpdateAvatar(ctx context.Context, id string, avatar string) (*mongo.UpdateResult, error)
	ChangeOwner(ctx context.Context, id string, owner entity.PersonInfo) (*mongo.UpdateResult, error)
	UpdateSetting(ctx context.Context, id string, setting entity.GroupSetting) (*mongo.UpdateResult, error)
	UpdateAvatarInOwner(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInOwner(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error)
	UpdateAvatarInAdmins(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInAdmins(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error)
	UpdateAvatarInMembers(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error)
	UpdateNameInMembers(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error)
	DeleteGroupByID(ctx context.Context, id string) error
}

type IVotingRepo interface {
	ExecTransaction(ctx context.Context, fn func(ctx context.Context) (any, error)) (any, error)
	SaveVoting(ctx context.Context, voting *entity.Voting) error
	FindVotingByID(ctx context.Context, id string) (*entity.Voting, error)
	AppendVoter(ctx context.Context, id string, name string, info entity.PersonInfo) (*mongo.UpdateResult, error)
	RemoveVoter(ctx context.Context, id string, name string, userID string) (*mongo.UpdateResult, error)
	LockVoting(ctx context.Context, id string, isLock bool, dateLock time.Time) (*mongo.UpdateResult, error)
	DeleteVotingByID(ctx context.Context, id string) error
}

type IUserSocketSvc interface {
	AppendFriendRequests(ctx context.Context, req dto.FriendRequestAddDTO) error
	RemoveFriendRequests(ctx context.Context, req dto.FriendRequestRemoveDTO) error
	AcceptFriendRequests(ctx context.Context, req dto.FriendRequestAcceptDTO) error
	Unfriend(ctx context.Context, req dto.UnfriendDTO) error
	UpdateTypeConversation(ctx context.Context, senderID, receiverID, typeSender, typeReceiver string) error
	AppendConversations(ctx context.Context, req dto.AppendConversationDTO, _type string) error
	AppendConversation(ctx context.Context, userID string, conversation entity.Conversation) error
	AppendConversationToMultiple(ctx context.Context, userID []string, conversation entity.Conversation) error
	RemoveConversation(ctx context.Context, userID string, idChat string) error
	RemoveConversationFromMultiple(ctx context.Context, userID []string, idChat string) error
	UpdateConversations(ctx context.Context, chat entity.Chat) error
	UpdateChatNameInConversation(ctx context.Context, arrID []string, chatID, chatName string) error
	UpdateAvatarInConversation(ctx context.Context, arrID []string, chatID, newAvatar string) error
}

type IChatSocketSvc interface {
	Create(ctx context.Context, chatID string) error
	Delete(ctx context.Context, chatID string) error
	AppendChat(ctx context.Context, chatID string, req dto.MessageAppendDTO) error
	ChangeDeliveryChat(ctx context.Context, chatID string, req dto.MessageDeliveryDTO) error
	ChangeReadChat(ctx context.Context, chatID string, req dto.MessageDeliveryDTO) error
	AppendHiddenMessage(ctx context.Context, chatID string, req dto.MessageHiddenDTO) error
	RecallMessage(ctx context.Context, chatID string, req dto.MessageHiddenDTO) error
	GetChatTop10(ctx context.Context, chatID string) (*entity.Chat, error)
	AppendVoter(ctx context.Context, req dto.AppendVoterDTO, chatID string, obj dto.MessageAppendDTO) error
	ChangeVoting(ctx context.Context, req dto.ChangeVoterDTO, chatID string, obj dto.MessageAppendDTO) error
	LockVoting(ctx context.Context, chatID string, req dto.MessageAppendDTO) error
}
type IGroupSocketSvc interface {
	Create(ctx context.Context, arrayID []string, req dto.CreateGroupDTO) ([]string, error)
	Delete(ctx context.Context, idChat string) ([]string, error)
	AppendMember(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error)
	AppendAdmin(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error)
	ChangeOwner(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error)
	RemoveMember(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error)
	RemoveAdmin(ctx context.Context, req dto.AppendMemberGroupDTO) ([]string, error)
	UpdateNameChat(ctx context.Context, req dto.ChangeNameChatGroupDTO) ([]string, error)
	UpdateAvatar(ctx context.Context, req dto.ChangeAvatarGroupDTO) ([]string, error)
	UpdateSettingChangeChatNameAndAvatar(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error)
	UpdateSettingPinMessages(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error)
	UpdateSettingSendMessages(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error)
	UpdateSettingMembershipApproval(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error)
	UpdateSettingCreateNewPolls(ctx context.Context, req dto.UpdateSettingGroupDTO) ([]string, error)
}

type IUpdateAsyncSvc interface {
	UpdateAvatarAsync(ctx context.Context, oldAvatar, newAvatar string)
	UpdateNameAsync(ctx context.Context, userID, newName string)
}
