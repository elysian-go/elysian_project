package project

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectRepository struct {
	DB *mongo.Database
}

func ProvideProjectRepository(db *mongo.Database) ProjectRepository {
	return ProjectRepository{DB: db}
}

func (r *ProjectRepository) FindAll() []Project {
	var projects []Project
	return projects
}

func (r *ProjectRepository) FindByID(id string) (Project, error) {
	return Project{}, nil
}

func (r *ProjectRepository) Save(project Project) (Project, error) {
	return Project{}, nil
}

func (r *ProjectRepository) Delete(project Project) {
}
