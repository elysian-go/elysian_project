//+build wireinject

package main

import (
	"github.com/VictorDebray/elysian_project/project"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/jinzhu/gorm"
)

func InitProjectAPI(rdb *gorm.DB, db *mongo.Database) project.ProjectAPI {
	panic(wire.Build(
		project.ProvideProjectRepository, project.ProvideProjectService, project.ProvideProjectAPI))
}