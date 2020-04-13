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
	projects := p.ProjectService.FindAll()

	c.JSON(http.StatusOK, gin.H{"projects": ToProjectModels(projects)})
}

func (p *ProjectAPI) FindByID(c *gin.Context) {
	id :=  c.Param("id")
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
