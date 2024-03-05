package router

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	_notesController "github.com/elizandrodantas/machine-controller-v2/notes/http/controller"
	_notesRepository "github.com/elizandrodantas/machine-controller-v2/notes/repository"
	_notesUsecase "github.com/elizandrodantas/machine-controller-v2/notes/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewNotesRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	notesRepo := _notesRepository.NewNotesRepository(pool)
	notesUC := _notesUsecase.NewNotesRepository(notesRepo)

	controller := _notesController.NotesController{
		NotesUsecase: notesUC,
		Timeout:      timeout,
	}

	g.POST("/notes", middleware.EnsureAutenticate(env), controller.Create)
	g.GET("/notes/:id", middleware.EnsureAutenticate(env), controller.List)
	g.PUT("/notes", middleware.EnsureAutenticate(env), controller.Update)
	g.DELETE("/notes", middleware.EnsureAutenticate(env), controller.Delete)
}
