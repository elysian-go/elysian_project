package project

func ToProject(projectModel Model) Project {
	return Project{
		Title:       projectModel.Title,
		Description: projectModel.Description,
		Owner:       projectModel.Owner,
		Archived:    projectModel.Archived,
	}
}

func ToProjectModel(project Project) Model {
	return Model{
		ID: project.ID,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}
}

func ToProjectModels(projects []Project) []Model {
	projectModels := make([]Model, len(projects))

	for i, itm := range projects {
		projectModels[i] = ToProjectModel(itm)
	}

	return projectModels
}