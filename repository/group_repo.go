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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repo) groupColl() *mongo.Collection {
	return r.db.Database(config.Cfg.DB.DBName).Collection("groups")
}

func (r *Repo) SaveGroup(ctx context.Context, group *entity.Group) (err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	defer errors.WrapDatabaseError(&err)

	opts := options.Update().SetUpsert(true)
	update := bson.D{
		{"$set", group},
	}
	_, err = r.groupColl().UpdateOne(ctx, bson.M{"id": group.ID}, update, opts)
	if err != nil {
		if strings.Contains(err.Error(), "E11000 duplicate key error collection") {
			return errors.GroupExists()
		}
		return err
	}
	return nil
}

func (r *Repo) FindGroupByID(ctx context.Context, id string) (*entity.Group, error) {
	var d entity.Group
	filter := bson.D{
		{"id", id},
	}
	if err := r.groupColl().FindOne(ctx, filter).Decode(&d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.GroupNotFound()
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repo) AppendMember(ctx context.Context, id string, info entity.PersonInfo) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$push": bson.M{"members": info}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveMember(ctx context.Context, id string, userID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$pull": bson.M{"members": bson.M{"user_id": userID}}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) AppendAdmin(ctx context.Context, id string, info entity.PersonInfo) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$push": bson.M{"admin": info}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) RemoveAdmin(ctx context.Context, id string, userID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$pull": bson.M{"admin": bson.M{"user_id": userID}}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateNameChat(ctx context.Context, id string, chatName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"chat_name": chatName}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatar(ctx context.Context, id string, avatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"chat_avatar": avatar}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) ChangeOwner(ctx context.Context, id string, owner entity.PersonInfo) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"owner": owner}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateSetting(ctx context.Context, id string, setting entity.GroupSetting) (*mongo.UpdateResult, error) {
	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"setting": setting}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInOwner(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"owner.user_id": userID}
	update := bson.M{"$set": bson.M{"owner.user_avatar": newAvatar}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateNameInOwner(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"owner.user_id": userID}
	update := bson.M{"$set": bson.M{"owner.user_name": newName}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInAdmins(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"admin.user_id": userID}
	update := bson.M{"$set": bson.M{"admin.$.user_avatar": newAvatar}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateNameInAdmins(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"admin.user_id": userID}
	update := bson.M{"$set": bson.M{"admin.$.user_name": newName}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateAvatarInMembers(ctx context.Context, userID string, newAvatar string) (*mongo.UpdateResult, error) {
	filter := bson.M{"members.user_id": userID}
	update := bson.M{"$set": bson.M{"members.$.user_avatar": newAvatar}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) UpdateNameInMembers(ctx context.Context, userID string, newName string) (*mongo.UpdateResult, error) {
	filter := bson.M{"members.user_id": userID}
	update := bson.M{"$set": bson.M{"members.$.user_name": newName}}
	return r.groupColl().UpdateOne(ctx, filter, update)
}

func (r *Repo) DeleteGroupByID(ctx context.Context, id string) error {
	filter := bson.D{{"id", id}}
	_, err := r.groupColl().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
