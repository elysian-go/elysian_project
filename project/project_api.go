package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type ProjectAPI struct {
	ProjectService ProjectService
}

func ProvideProjectAPI(p ProjectService) ProjectAPI {
	return ProjectAPI{ProjectService: p}
}

func (p *ProjectAPI) FindAll(c *gin.Context) {
	value := c.MustGet("user_id")
	userId, ok := value.(string)
	if !ok {
		log.Printf("got data of type %T but wanted string", value)
		c.JSON(http.StatusInternalServerError, "internal error")
		return
	}
	var projectsOwned, projectsCollaborator []Project
	var err error
	projectsOwned, err = p.ProjectService.FindProjectsByOwner(userId)
	projectsCollaborator, err = p.ProjectService.FindProjectsByCollaborator(userId)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projectsOwned": ToProjectModels(projectsOwned),
		"projectsCollaborator": ToProjectModels(projectsCollaborator)})
}

func (p *ProjectAPI) FindByID(c *gin.Context) {
	id :=  c.Param("pid")
	project, err := p.ProjectService.FindByID(id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": ToProjectModel(project)})
}

func (p *ProjectAPI) Create(c *gin.Context) {
	var projectModel Model
	err := c.BindJSON(&projectModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//Getting user id (aka owner id) from context
	value := c.MustGet("user_id")
	userId, ok := value.(string)
	if !ok {
		log.Printf("got data of type %T but wanted string", value)
		c.JSON(http.StatusInternalServerError, "internal error")
		return
	}
	projectModel.Owner = userId
	project, err := p.ProjectService.Save(ToProject(projectModel))
	if err != nil {
		log.Println(err)
		switch {
		case strings.Contains(err.Error(), "duplicate"):
			c.JSON(http.StatusConflict, gin.H{"error": "project has already an owner"})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (p *ProjectAPI) AddCollaborator(c *gin.Context) {
	var addContributorModel AddContributorModel
	err := c.BindJSON(&addContributorModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projectId := c.Param("pid")
	value := c.MustGet("user_id")
	userId, ok := value.(string)
	if !ok {
		log.Printf("got data of type %T but wanted string", value)
		c.JSON(http.StatusInternalServerError, "internal error")
		return
	}

	projectOwner, err := p.ProjectService.FindOwnerProjectByIds(userId, projectId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access forbidden, please make sure you are logged in" +
			" and the owner of the project"})
		return
	}

	err = p.ProjectService.SaveCollaborator(projectOwner.ProjectId, addContributorModel.ID[0])
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate"):
			c.JSON(http.StatusConflict, gin.H{"error": "user is already contributor of project"})
		default:
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "contributors added successfully"})
}

func (p *ProjectAPI) Delete(c *gin.Context) {
	//id := c.Param("id")
	project := Project{} //p.ProjectService.FindByID(id)

	if project == (Project{}) {
		c.Status(http.StatusBadRequest)
		return
	}

	p.ProjectService.Delete(project)

	c.Status(http.StatusOK)
}
