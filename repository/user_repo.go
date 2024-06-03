package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) userColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("users")
}

// func (r *Repo) CreateUserIndexes(ctx context.Context) ([]string, error) {
// 	indexes, err := r.userColl().Indexes().CreateMany(ctx, []mongo.IndexModel{
// 		{Keys: bson.D{
// 			{"id", 1},
// 			{"deleted_at", 1},
// 		}, Options: options.Index().SetUnique(true)},
// 		{Keys: bson.D{
// 			{"username", 1},
// 			{"deleted_at", 1},
// 		}, Options: options.Index().SetUnique(true)},
// 		{Keys: bson.D{
// 			{"email", 1},
// 			{"deleted_at", 1},
// 		}, Options: options.Index().SetUnique(true)},
// 		{Keys: bson.D{
// 			{"name", "text"},
// 			{"username", "text"},
// 			{"email", "text"},
// 		}},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return indexes, nil
// }

func (r *Repo) SaveUser(ctx context.Context, user *entity.User) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", user},
	}
	_, err = r.userColl().UpdateOne(ctx, bson.M{"id": user.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.UserExists()
		}
		return err
	}
	return nil
}

func (r *Repo) FindUserByID(ctx context.Context, id string) (*entity.User, error) {
	filter := bson.M{"id": id}
	findOptions := options.FindOne().SetSort(bson.D{{"conversations.updated_at", -1}})
	var user entity.User
	err := r.userColl().FindOne(ctx, filter, findOptions).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.UserNotFound()
		}
		return nil, err
	}
	return &user, err
}

func (r *Repo) AppendFriendRequest(ctx context.Context, id string, friendRequest entity.FriendRequest) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$push": bson.M{"friend_requests": friendRequest}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveFriendRequest(ctx context.Context, senderID, receiver string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": senderID}
	update := bson.M{"$pull": bson.M{"friend_requests": bson.M{"user_id": receiver}}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateTypeConversation(ctx context.Context, senderID, chatID1, chatID2, convoType string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": senderID, "conversations.chat_id": bson.M{"$in": []string{chatID1, chatID2}}}
	update := bson.M{"$set": bson.M{"conversations.$.type": convoType}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) AppendConversation(ctx context.Context, id string, conversation entity.Conversation) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$push": bson.M{"conversations": conversation}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) AppendConversationToMultiple(ctx context.Context, ids []string, conversation entity.Conversation) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": bson.M{"$in": ids}}
	update := bson.M{"$push": bson.M{"conversations": conversation}}
	return r.userColl().UpdateMany(ctx, filter, update)
}

func (r *Repo) RemoveConversation(ctx context.Context, id, chatID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$pull": bson.M{"conversations": bson.M{"chat_id": chatID}}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveConversationFromMultiple(ctx context.Context, ids []string, chatID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": bson.M{"$in": ids}}
	update := bson.M{"$pull": bson.M{"conversations": bson.M{"chat_id": chatID}}}
	return r.userColl().UpdateMany(ctx, filter, update)
}

func (r *Repo) SearchConversation(ctx context.Context, senderID, chatID1, chatID2 string) (*entity.User, error) {
	filter := bson.M{"id": senderID, "conversations.chat_id": bson.M{"$in": []string{chatID1, chatID2}}}
	var user entity.User
	err := r.userColl().FindOne(ctx, filter).Decode(&user)
	return &user, err
}

func (r *Repo) SearchSingleConversation(ctx context.Context, senderID, chatID string) (*entity.User, error) {
	filter := bson.M{"id": senderID, "conversations.chat_id": chatID}
	var user entity.User
	err := r.userColl().FindOne(ctx, filter).Decode(&user)
	return &user, err
}

func (r *Repo) UpdateChatActivity(ctx context.Context, chatID string, lastUpdateAt time.Time, deliveries, reads []entity.Delivery, newTopChatActivity []entity.ChatActivity) (*mongo.UpdateResult, error) {
	filter := bson.M{"conversations.chat_id": chatID}
	update := bson.M{
		"$set": bson.M{
			"conversations.$.updated_at":          lastUpdateAt,
			"conversations.$.deliveries":          deliveries,
			"conversations.$.reads":               reads,
			"conversations.$.top_chat_activities": newTopChatActivity,
		},
	}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInConversation(ctx context.Context, userID, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"conversations.id_user_or_group": userID}
	update := bson.M{"$set": bson.M{"conversations.$.chat_avatar": newAvatar}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateNameInConversation(ctx context.Context, userID, newName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"conversations.id_user_or_group": userID}
	update := bson.M{"$set": bson.M{"conversations.$.chat_name": newName}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInFriendRequest(ctx context.Context, userID, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"friendRequests.user_id": userID}
	update := bson.M{"$set": bson.M{"friendRequests.$.user_avatar": newAvatar}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateNameInFriendRequest(ctx context.Context, userID, newName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"friendRequests.user_id": userID}
	update := bson.M{"$set": bson.M{"friendRequests.$.user_name": newName}}
	return r.userColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateChatNameInConversation(ctx context.Context, ids []string, chatID, chatName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": bson.M{"$in": ids}, "conversations.chat_id": chatID}
	update := bson.M{"$set": bson.M{"conversations.$.chat_name": chatName}}
	return r.userColl().UpdateMany(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInConversationMultiple(ctx context.Context, ids []string, chatID, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": bson.M{"$in": ids}, "conversations.chat_id": chatID}
	update := bson.M{"$set": bson.M{"conversations.$.chat_avatar": newAvatar}}
	return r.userColl().UpdateMany(ctx, filter, update)
}

func (r *Repo) DeleteUserByID(ctx context.Context, id string) error {
	filter := bson.D{{"id", id}}
	_, err := r.userColl().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAllRecentSearchProfiles(ctx context.Context, userID string) ([]entity.Profile, error) {
	filter := bson.M{"id": userID}
	var user entity.User
	err := r.userColl().FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.UserNotFound()
		}
		return nil, err
	}
	return user.RecentSearchProfiles, err
}

func (r *Repo) UpdateRecentSearchProfiles(ctx context.Context, userID string, recentSearchProfiles []entity.Profile) error {
	filter := bson.M{"id": userID}
	update := bson.M{"$set": bson.M{"recent_search_profiles": recentSearchProfiles}}
	_, err := r.userColl().UpdateOne(ctx, filter, update)
	return err
}
