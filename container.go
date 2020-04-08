//+build wireinject

package main

import (
	"github.com/VictorDebray/elysian_project/project"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitProjectAPI(db *mongo.Database) project.ProjectAPI {
	panic(wire.Build(
		project.ProvideProjectRepository, project.ProvideProjectService, project.ProvideProjectAPI))
}