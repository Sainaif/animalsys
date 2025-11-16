package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sainaif/animalsys/backend/internal/domain/entities"
	animalUC "github.com/sainaif/animalsys/backend/internal/usecase/animal"
	dashboardUC "github.com/sainaif/animalsys/backend/internal/usecase/dashboard"
	taskUC "github.com/sainaif/animalsys/backend/internal/usecase/task"
)

type DashboardHandler struct {
	dashboardUseCase *dashboardUC.DashboardUseCase
	animalUseCase    *animalUC.AnimalUseCase
	taskUseCase      *taskUC.TaskUseCase
}

type dashboardAnimalSummary struct {
	ID          string                    `json:"id"`
	Name        entities.MultilingualName `json:"name"`
	Species     string                    `json:"species"`
	Breed       string                    `json:"breed,omitempty"`
	Status      entities.AnimalStatus     `json:"status"`
	Color       string                    `json:"color,omitempty"`
	Sex         entities.AnimalSex        `json:"sex"`
	DateOfBirth *time.Time                `json:"date_of_birth,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
	Images      entities.AnimalImages     `json:"images"`
}

type dashboardTaskSummary struct {
	ID         string                `json:"id"`
	Title      string                `json:"title"`
	Category   entities.TaskCategory `json:"category"`
	Status     entities.TaskStatus   `json:"status"`
	Priority   entities.TaskPriority `json:"priority"`
	DueDate    *time.Time            `json:"due_date,omitempty"`
	AssignedTo string                `json:"assigned_to,omitempty"`
}

func NewDashboardHandler(
	dashboardUseCase *dashboardUC.DashboardUseCase,
	animalUseCase *animalUC.AnimalUseCase,
	taskUseCase *taskUC.TaskUseCase,
) *DashboardHandler {
	return &DashboardHandler{
		dashboardUseCase: dashboardUseCase,
		animalUseCase:    animalUseCase,
		taskUseCase:      taskUseCase,
	}
}

// GetDashboardMetrics gets full dashboard metrics
func (h *DashboardHandler) GetDashboardMetrics(c *gin.Context) {
	metrics, err := h.dashboardUseCase.GetDashboardMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard metrics"})
		return
	}

	recentAnimals, err := h.getRecentAnimals(c.Request.Context())
	if err != nil {
		log.Printf("failed to load recent animals for dashboard: %v", err)
	}

	upcomingTasks, err := h.getUpcomingTasks(c.Request.Context())
	if err != nil {
		log.Printf("failed to load upcoming tasks for dashboard: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"overview":        metrics.Overview,
		"animals":         metrics.Animals,
		"adoptions":       metrics.Adoptions,
		"donations":       metrics.Donations,
		"volunteers":      metrics.Volunteers,
		"events":          metrics.Events,
		"veterinary":      metrics.Veterinary,
		"recent_activity": metrics.RecentActivity,
		"recent_animals":  recentAnimals,
		"upcoming_tasks":  upcomingTasks,
		"trends":          metrics.Trends,
		"generated_at":    metrics.GeneratedAt,
	})
}

// GetOverviewMetrics gets overview metrics only
func (h *DashboardHandler) GetOverviewMetrics(c *gin.Context) {
	overview, err := h.dashboardUseCase.GetOverviewMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load overview metrics"})
		return
	}

	c.JSON(http.StatusOK, overview)
}

// GetDashboardSummary gets dashboard summary
func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	metrics, err := h.dashboardUseCase.GetDashboardMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"overview":   metrics.Overview,
		"animals":    metrics.Animals,
		"adoptions":  metrics.Adoptions,
		"donations":  metrics.Donations,
		"volunteers": metrics.Volunteers,
	})
}

// GetDashboardWidgets gets dashboard widgets
func (h *DashboardHandler) GetDashboardWidgets(c *gin.Context) {
	metrics, err := h.dashboardUseCase.GetDashboardMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load dashboard widgets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recent_activity": metrics.RecentActivity,
		"trends":          metrics.Trends,
	})
}

func (h *DashboardHandler) getRecentAnimals(ctx context.Context) ([]dashboardAnimalSummary, error) {
	if h.animalUseCase == nil {
		return []dashboardAnimalSummary{}, nil
	}

	resp, err := h.animalUseCase.ListAnimals(ctx, &animalUC.ListAnimalsRequest{
		Limit:     5,
		SortBy:    "created_at",
		SortOrder: "desc",
	})
	if err != nil {
		return nil, err
	}

	items := make([]dashboardAnimalSummary, 0, len(resp.Animals))
	for _, animal := range resp.Animals {
		items = append(items, dashboardAnimalSummary{
			ID:          animal.ID.Hex(),
			Name:        animal.Name,
			Species:     animal.Species,
			Breed:       animal.Breed,
			Status:      animal.Status,
			Color:       animal.Color,
			Sex:         animal.Sex,
			DateOfBirth: animal.DateOfBirth,
			CreatedAt:   animal.CreatedAt,
			Images:      animal.Images,
		})
	}

	return items, nil
}

func (h *DashboardHandler) getUpcomingTasks(ctx context.Context) ([]dashboardTaskSummary, error) {
	if h.taskUseCase == nil {
		return []dashboardTaskSummary{}, nil
	}

	tasks, err := h.taskUseCase.GetUpcomingTasks(ctx, 14)
	if err != nil {
		return nil, err
	}

	items := make([]dashboardTaskSummary, 0, len(tasks))
	for _, task := range tasks {
		taskSummary := dashboardTaskSummary{
			ID:       task.ID.Hex(),
			Title:    task.Title,
			Category: task.Category,
			Status:   task.Status,
			Priority: task.Priority,
			DueDate:  task.DueDate,
		}

		if task.AssignedTo != nil {
			taskSummary.AssignedTo = task.AssignedTo.Hex()
		}

		items = append(items, taskSummary)
	}

	return items, nil
}
