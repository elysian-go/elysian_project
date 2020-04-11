// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/VictorDebray/elysian_project/project"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from container.go:

func InitProjectAPI(db *mongo.Database) project.ProjectAPI {
	projectRepository := project.ProvideProjectRepository(db)
	projectService := project.ProvideProjectService(projectRepository)
	projectAPI := project.ProvideProjectAPI(projectService)
	return projectAPI
}