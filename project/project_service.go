package project

import (
	"errors"
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
		return project, errors.New("resource not found")
	}
	return project, err
}

func (s *ProjectService) Save(project Project) (Project, error) {
	project, err := s.ProjectRepository.Save(project)
	if err != nil {
		return project, errors.New("duplicate entry on email")
	}
	return project, nil
}

func (s *ProjectService) Update(project Project) (Project, error) {
	return project, nil
}

func (s *ProjectService) Delete(project Project) {
	s.ProjectRepository.Delete(project)
}
