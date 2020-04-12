package project

import (
	"errors"
	"log"
)

type ProjectService struct {
	ProjectRepository ProjectRepository
}

func ProvideProjectService(p ProjectRepository) ProjectService {
	return ProjectService{ProjectRepository: p}
}

func (s *ProjectService) FindAll() []Project {
	return s.ProjectRepository.FindAll()
}

func (s *ProjectService) FindByID(id string) (Project, error) {
	project, err := s.ProjectRepository.FindByID(id)
	if err != nil {
		log.Println(err.Error())
		return project, errors.New("resource not found")
	}
	return project, err
}

func (s *ProjectService) Save(project Project) (Project, error) {
	projectId, err := s.ProjectRepository.Save(project)
	if err != nil {
		return Project{}, errors.New("could not save project to database")
	}
	log.Println(projectId)
	projectSaved, err := s.FindByID(projectId)
	if err != nil {
		return Project{}, errors.New("could not retrieve saved document")
	}
	return projectSaved, nil
}

func (s *ProjectService) Update(project Project) (Project, error) {
	return project, nil
}

func (s *ProjectService) Delete(project Project) {
	s.ProjectRepository.Delete(project)
}
