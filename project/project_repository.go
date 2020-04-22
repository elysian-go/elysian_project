package project

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *ProjectRepository) FindProjectsByOwnerId(accountId string) ([]Project, error) {
	var projects []Project
	rows, err := r.RDB.Raw("SELECT project_id FROM owner_project WHERE user_id = ?", accountId).Rows()
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "error in find project owner rows")
	}

	for rows.Next() {
		var projectId string
		rows.Scan(&projectId)
		project, err := r.FindByID(projectId)
		if err != nil {
			return nil, errors.Wrap(err, "error while retrieving project by id from row")
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepository) FindProjectsByCollaboratorId(accountId string) ([]Project, error) {
	var projects []Project
	rows, err := r.RDB.Raw("SELECT project_id FROM collaborator_project WHERE user_id = ?", accountId).Rows()
	defer rows.Close()
	if err != nil {
		return nil, errors.Wrap(err, "error in find project owner rows")
	}

	for rows.Next() {
		var projectId string
		rows.Scan(&projectId)
		project, err := r.FindByID(projectId)
		if err != nil {
			return nil, errors.Wrap(err, "error while retrieving project by id from row")
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepository) FindOwnerProjectRelation(ownerId string, projectId string) (OwnerProject, error) {
	query := r.RDB.Raw("SELECT * FROM owner_project WHERE user_id = ? AND project_id = ?", ownerId, projectId)
	err := query.Error
	if err != nil {
		return OwnerProject{}, errors.Wrap(err, "find OwnerProject by ids failed")
	}

	var ownerProject OwnerProject
	err = query.Scan(&ownerProject).Error
	if err != nil {
		return OwnerProject{}, errors.Wrap(err, "scan query to owner project failed")
	}
	return ownerProject, nil
}

func (r *ProjectRepository) FindByID(id string) (Project, error) {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objId}}

	var project Project
	res := r.collection.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return Project{}, errors.Wrap(res.Err(), "project not found by id")
	}

	err := res.Decode(&project)
	if err != nil {
		return Project{}, errors.Wrap(err,"error while decoding result to project obj")
	}
	return project, nil
}

func (r *ProjectRepository) SaveProject(project Project) (string, error) {
	result, err := r.collection.InsertOne(context.TODO(), project)
	if err != nil {
		return "", errors.Wrap(err, "error while inserting project")
	}

	var resId, ok = result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.Wrap(err, "error while casting InsertedID to ObjectID")
	}
	return resId.Hex(), nil
}

func (r *ProjectRepository) SaveOwnerProject(ownerProject OwnerProject) (OwnerProject, error) {
	err := r.RDB.Save(&ownerProject).Error
	if err != nil {
		return OwnerProject{}, errors.Wrap(err, "error while saving owner project")
	}
	return ownerProject, nil
}

func (r *ProjectRepository) SaveCollaborator(projectId string, collaboratorId string) error {
	query := r.RDB.Exec("INSERT INTO collaborator_project (user_id, project_id) VALUES (?, ?)",
		collaboratorId, projectId)
	err := query.Error
	if err != nil {
		return errors.Wrap(err, "error while saving collaborator")
	}
	return nil
}

func (r *ProjectRepository) Save(project Project) (string, error) {
	projectId, err := r.SaveProject(project)
	if err != nil {
		return "", errors.Wrap(err, "error while inserting project in transaction")
	}

	ownerProj := OwnerProject{
		UserId:    project.Owner,
		ProjectId: projectId,
	}
	if err = r.RDB.Create(&ownerProj).Error; err != nil {
		r.Delete(projectId)
		return "", errors.Wrap(err, "error while creating OwnerProject in transaction")
	}
	return projectId, err
}

func (r *ProjectRepository) Delete(id string) error {
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", objId}}
	_, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Wrap(err, "error while deleting project")
	}
	return nil
}
