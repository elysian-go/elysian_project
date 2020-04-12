package project

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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

func (r *ProjectRepository) FindByProjectsAccountId(accountId string) ([]Project, error) {
	var projects []Project
	rows, err := r.RDB.Raw("SELECT * FROM account WHERE id = ?", accountId).Rows()
	defer rows.Close()
	if err != nil {
		return nil, errors.New("internal error")
	}

	for rows.Next() {
		var ownerProj OwnerProject
		rows.Scan(&ownerProj)
		project, err := r.FindByID(ownerProj.ProjectId)
		if err != nil {
			return nil, errors.New("error while retrieving project by id")
		}
		projects = append(projects, project)
	}
	return projects, nil
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

func (r *ProjectRepository) SaveProject(project Project) (string, error) {
	result, err := r.collection.InsertOne(context.TODO(), project)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("error while inserting project")
	}

	var resId, ok = result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("error while casting InsertedID to ObjectID")
	}
	return resId.Hex(), nil
}

func (r *ProjectRepository) SaveOwnerProject(ownerProject OwnerProject) (OwnerProject, error) {
	err := r.RDB.Save(&ownerProject).Error
	if err != nil {
		log.Println(err.Error())
		return OwnerProject{}, err
	}
	return ownerProject, nil
}

func (r *ProjectRepository) Save(project Project) (string, error) {
	projectId, err := r.SaveProject(project)
	if err != nil {
		log.Println(err.Error())
		log.Println("error while inserting project in transaction")
		return "", err
	}

	ownerProj := OwnerProject{
		UserId:    project.Owner,
		ProjectId: projectId,
	}
	if err = r.RDB.Create(&ownerProj).Error; err != nil {
		log.Println(err.Error())
		log.Println("error while creating OwnerProject in transaction")
		r.Delete(projectId)
		return "", err
	}
	return projectId, err
}

func (r *ProjectRepository) Delete(id string) error {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objId}}
	_, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
