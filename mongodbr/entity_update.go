package mongodbr

import (
	"github.com/shanluzhineng/fwpkg/mongodbr/builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// update
type IEntityUpdate interface {
	FindOneAndUpdate(entity IEntity, opts ...*options.FindOneAndUpdateOptions) error
	FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error
	UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (interface{}, error)
}

var _ IEntityUpdate = (*MongoCol)(nil)

// #region update members

func (r *MongoCol) FindOneAndUpdate(entity IEntity, opts ...*options.FindOneAndUpdateOptions) error {
	// if entity == nil {
	// 	return fmt.Errorf("在更新%s数据时item参数不能为nil", r.documentName)
	// }

	objectId := entity.GetObjectId()
	update := builder.NewBsonBuilder().NewOrUpdateSet(entity).ToValue()
	return r.FindOneAndUpdateWithId(objectId, update, opts...)
}

func (r *MongoCol) FindOneAndUpdateWithId(objectId primitive.ObjectID, update interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	//没有设置参数，使用默认的
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	if len(opts) <= 0 {
		opts = make([]*options.FindOneAndUpdateOptions, 0)
		opts = append(opts, options.FindOneAndUpdate().SetUpsert(false))
	}
	// 检查传入数据，只更新非 _id 字段
	upMap, err := ToMap(update)
	if err != nil {
		return err
	}
	if _, ok := upMap["$set"]; !ok {
		_, incOk := upMap["$inc"]
		_, pushOk := upMap["$push"]
		_, pullOk := upMap["$pull"]
		_, popOk := upMap["$pop"]
		_, unsetOk := upMap["$unset"]
		if !incOk && !pushOk && !pullOk && !popOk && !unsetOk {
			delete(upMap, "_id")
			upMap = bson.M{"$set": upMap} // 无其他特殊操作符时再加set。。
		}
	}
	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectId},
		upMap,
		opts...,
	).Err(); err != nil {
		return err
	}
	return nil
}

func (r *MongoCol) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoCol) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (interface{}, error) {
	ctx, cancel := CreateContext(r.configuration)
	defer cancel()

	result, err := r.collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		if result != nil {
			return result.UpsertedID, err
		} else {
			return nil, err
		}
	}

	if result != nil {
		return result.UpsertedID, nil
	}
	return nil, nil
}

// #endregion

func ToMap(data interface{}) (map[string]interface{}, error) {
	var m = make(map[string]interface{})
	bt, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := bson.Unmarshal(bt, &m); err != nil {
		return nil, err
	}
	return m, nil
}
