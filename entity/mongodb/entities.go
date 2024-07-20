package mongodb

import (
	"github.com/shanluzhineng/fwpkg/mongodbr"
	"github.com/shanluzhineng/fwpkg/system/reflector"
	"github.com/shanluzhineng/fwpkg/utils/str"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollectionName(v interface{}) string {
	return str.ToSnake(reflector.GetName(v))
}

// regist Repository create option
func RegistEntityRepositoryOption[T mongodbr.IEntity](clientKey string, databaseName string, opts ...mongodbr.RepositoryOption) {
	d := GetDatabase(clientKey, databaseName)
	collectionName := GetCollectionName((new(T)))
	d._entityRepositoryOptionMap[collectionName] = createEntityRepositoryOption[T](opts...)
	d._repositoryMapping[collectionName] = d.ensureCreateRepository(collectionName, opts...)
}

func createEntityRepositoryOption[T mongodbr.IEntity](opts ...mongodbr.RepositoryOption) []mongodbr.RepositoryOption {
	if len(opts) <= 0 {
		opts = append(opts, mongodbr.WithCreateItemFunc(func() interface{} {
			return new(T)
		}))
		opts = append(opts, mongodbr.WithDefaultSort(func(fo *options.FindOptions) *options.FindOptions {
			return fo.SetSort(bson.D{{Key: "_id", Value: -1}})
		}))
	}
	return opts
}
