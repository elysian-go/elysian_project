package project

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectRepository struct {
	RDB        *gorm.DB
	DB         *mongo.Database
	collection *mongo.Collection
}

func ProvideProjectRepository(rdb *gorm.DB, db *mongo.Database) ProjectRepository {
	return ProjectRepository{RDB: rdb, DB: db, collection: db.Collection("project")}
}

func (r *ProjectRepository) FindAll() []Project {
	var projects []Project
	return projects
}

func (r *ProjectRepository) FindByID(id string) (Project, error) {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objId}}

	var project Project
	res := r.collection.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return Project{}, errors.New("project not found by id")
	}

	err := res.Decode(&project)
	if err != nil {
		return Project{}, errors.New("error while decoding result to project obj")
	}
	return project, nil
}

func (r *ProjectRepository) Save(project Project) (string, error) {
	result, err := r.collection.InsertOne(context.TODO(), project)
	if err != nil {
		return "", errors.New("internal error")
	}

	var resId, ok = result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("internal error")
	}
	return resId.Hex(), nil
}

func (r *ProjectRepository) Delete(project Project) {
}
