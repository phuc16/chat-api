package repository

import (
	"app/config"
	"app/entity"
	"app/errors"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchIndexes struct {
	Index int `json:"index"`
}

func (r *Repo) chatColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("chats")
}

func (r *Repo) SaveChat(ctx context.Context, chat *entity.Chat) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", chat},
	}
	_, err = r.chatColl().UpdateOne(ctx, bson.M{"id": chat.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.ChatExists()
		}
		return err
	}
	return nil
}

func (r *Repo) FindChatByID(ctx context.Context, id string) (*entity.Chat, error) {
	var d entity.Chat
	filter := bson.D{
		{"id", id},
	}
	if err := r.chatColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ChatNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) DeleteChatByID(ctx context.Context, chatID string) error {
	filter := bson.D{{"id", chatID}}
	_, err := r.chatColl().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) AppendChatActivityByIDChat(ctx context.Context, chatID string, chatActivity entity.ChatActivity) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID}
	update := bson.M{"$push": bson.M{"chat_activities": chatActivity}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) AppendDelivery(ctx context.Context, chatID string, delivery entity.Delivery) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID}
	update := bson.M{"$push": bson.M{"deliveries": delivery}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) AppendRead(ctx context.Context, chatID string, delivery entity.Delivery) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID}
	update := bson.M{"$push": bson.M{"reads": delivery}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) SearchDeliveryByUserID(ctx context.Context, chatID string, userID string) (*entity.Chat, error) {
	filter := bson.M{"id": chatID, "deliveries.user_id": userID}
	var chat entity.Chat
	err := r.chatColl().FindOne(ctx, filter).Decode(&chat)
	return &chat, err
}

func (r *Repo) SearchReadByUserID(ctx context.Context, chatID string, userID string) (*entity.Chat, error) {
	filter := bson.M{"id": chatID, "reads.user_id": userID}
	var chat entity.Chat
	err := r.chatColl().FindOne(ctx, filter).Decode(&chat)
	return &chat, err
}

func (r *Repo) ChangeDelivery(ctx context.Context, chatID string, userID string, messageID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID, "deliveries.user_id": userID}
	update := bson.M{"$set": bson.M{"deliveries.$.message_id": messageID}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) ChangeRead(ctx context.Context, chatID string, userID string, messageID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID, "reads.user_id": userID}
	update := bson.M{"$set": bson.M{"reads.$.message_id": messageID}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveDelivery(ctx context.Context, chatID string, messageID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID}
	update := bson.M{"$pull": bson.M{"deliveries": bson.M{"message_id": messageID}}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveRead(ctx context.Context, chatID string, messageID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID}
	update := bson.M{"$pull": bson.M{"reads": bson.M{"message_id": messageID}}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) AppendHiddenMessage(ctx context.Context, chatID string, userID string, messageID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID, "chat_activities.message_id": messageID}
	update := bson.M{"$push": bson.M{"chat_activities.$.hidden": userID}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RecallMessage(ctx context.Context, chatID string, messageID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": chatID, "chat_activities.message_id": messageID}
	update := bson.M{"$set": bson.M{"chat_activities.$.recall": true}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) GetChatTop10(ctx context.Context, chatID string) (*entity.Chat, error) {
	filter := bson.M{"id": chatID}
	options := options.FindOne().SetProjection(bson.M{"chat_activities": bson.M{"$slice": -10}})
	var chat entity.Chat
	err := r.chatColl().FindOne(ctx, filter, options).Decode(&chat)
	return &chat, err
}

func (r *Repo) GetChatActivityFromNToM(ctx context.Context, chatID string, x int, y int) ([]entity.ChatActivity, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"id": chatID}}},
		{{"$project", bson.M{"id": 0, "chat_activities": 1}}},
		{{"$unwind", "$chat_activities"}},
		{{"$replaceRoot", bson.M{"newRoot": "$chat_activities"}}},
		{{"$sort", bson.M{"timestamp": -1}}},
		{{"$skip", x}},
		{{"$limit", y}},
		{{"$sort", bson.M{"timestamp": 1}}},
	}
	cursor, err := r.chatColl().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var activities []entity.ChatActivity
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *Repo) UpdateAvatarInRead(ctx context.Context, oldAvatar string, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"reads.user_avatar": oldAvatar}
	update := bson.M{"$set": bson.M{"reads.$.user_avatar": newAvatar}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInDelivery(ctx context.Context, oldAvatar string, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"deliveries.user_avatar": oldAvatar}
	update := bson.M{"$set": bson.M{"deliveries.$.user_avatar": newAvatar}}
	return r.chatColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) SearchByKeyWord(ctx context.Context, chatID string, key string) ([]entity.ChatActivity, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"id": chatID}}},
		{{"$unwind", "$chat_activities"}},
		{{"$replaceRoot", bson.M{"newRoot": "$chat_activities"}}},
		{{"$match", bson.M{"contents.value": primitive.Regex{Pattern: key, Options: "i"}}}},
		{{"$sort", bson.M{"timestamp": -1}}},
	}
	cursor, err := r.chatColl().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var activities []entity.ChatActivity
	if err = cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *Repo) GetIndexOfMessageID(ctx context.Context, chatID string, messageID string) ([]SearchIndexes, error) {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"id": chatID}}},
		{{"$unwind", "$chat_activities"}},
		{{"$replaceRoot", bson.M{"newRoot": "$chat_activities"}}},
		{{"$sort", bson.M{"timestamp": -1}}},
		{{"$group", bson.M{"id": nil, "chat_activities": bson.M{"$push": "$$ROOT"}}}},
		{{"$project", bson.M{"chat_activities": bson.M{"$map": bson.M{"input": "$chat_activities", "as": "activity", "in": bson.M{"$mergeObjects": bson.A{"$$activity", bson.M{"index": bson.M{"$indexOfArray": bson.A{"$chat_activities.message_id", "$$activity.message_id"}}}}}}}}}},
		{{"$unwind", "$chat_activities"}},
		{{"$replaceRoot", bson.M{"newRoot": "$chat_activities"}}},
		{{"$match", bson.M{"messageID": messageID}}},
		{{"$project", bson.M{"index": 1}}},
	}
	cursor, err := r.chatColl().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var result []SearchIndexes
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
