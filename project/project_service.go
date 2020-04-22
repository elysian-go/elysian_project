package project

import (
	"github.com/pkg/errors"
)

type ProjectService struct {
	ProjectRepository ProjectRepository
}

func ProvideProjectService(p ProjectRepository) ProjectService {
	return ProjectService{ProjectRepository: p}
}

func (s *ProjectService) FindByID(id string) (Project, error) {
	project, err := s.ProjectRepository.FindByID(id)
	if err != nil {
		return project, errors.Wrap(err, "could not find project by id")
	}
	return project, nil
}

// FindProjectsByCollaborator find collaborator projects of given user id
func (s *ProjectService) FindProjectsByCollaborator(ownerId string) ([]Project, error) {
	projects, err := s.ProjectRepository.FindProjectsByCollaboratorId(ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "could not find projects by owner")
	}
	return projects, nil
}

// FindProjectsByOwner find projects of given user id
func (s *ProjectService) FindProjectsByOwner(ownerId string) ([]Project, error) {
	projects, err := s.ProjectRepository.FindProjectsByOwnerId(ownerId)
	if err != nil {
		return nil, errors.Wrap(err, "could not find projects by owner")
	}
	return projects, nil
}

// FindOwnerProjectByIds find project owner relation by ids
func (s *ProjectService) FindOwnerProjectByIds(ownerId string, projectId string) (OwnerProject, error) {
	ownerProject, err := s.ProjectRepository.FindOwnerProjectRelation(ownerId, projectId)
	if err != nil {
		return OwnerProject{}, errors.Wrap(err, "could not find owner-project by ids")
	}
	return ownerProject, nil
}

// Save save project
func (s *ProjectService) Save(project Project) (Project, error) {
	projectId, err := s.ProjectRepository.Save(project)
	if err != nil {
		return Project{}, errors.Wrap(err, "could not save project to database")
	}
	projectSaved, err := s.FindByID(projectId)
	if err != nil {
		return Project{}, errors.Wrap(err, "could not retrieve saved document")
	}
	return projectSaved, nil
}

// SaveCollaborator save collaborator project relation
func (s *ProjectService) SaveCollaborator(projectId string, collaboratorId string) error {
	project, err := s.ProjectRepository.FindByID(projectId)
	if err != nil {
		return errors.Wrap(err, "could not find project by id")
	}
	if project.ID == "" {
		return nil
	}
	//TODO contributor user exists, with signed JWT ?
	err = s.ProjectRepository.SaveCollaborator(projectId, collaboratorId)
	if err != nil {
		return errors.Wrap(err, "could not save save collaborator")
	}
	return nil
}

func (s *ProjectService) Update(project Project) (Project, error) {
	return project, nil
}

func (s *ProjectService) Delete(project Project) {
	s.ProjectRepository.Delete(project.ID)
}
