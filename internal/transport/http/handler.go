package http

import (
	"encoding/json"
	"net/http"

	"calculator/internal/app/domain"
	"calculator/internal/app/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Calculator
}

func NewHandler(service *service.Calculator) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/calculate", h.Calculate).Methods("POST")
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

type Request struct {
	Instructions []domain.Instruction `json:"instructions"`
}

type Response struct {
	Items []domain.ResultItem `json:"items"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"ошибка"`
}

// Calculate обрабатывает HTTP-запрос для выполнения вычислений
// @Summary Выполнить вычисления
// @Description Принимает список инструкций и возвращает результаты вычислений
// @Tags Калькулятор
// @Accept json
// @Produce json
// @Param instructions body []domain.Instruction true "Список инструкций для выполнения"
// @Success 200 {object} Response "Успешный ответ с результатами"
// @Failure 400 {object} ErrorResponse "Ошибка в формате запроса"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /calculate [post]
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	results, err := h.service.ProcessInstructions(req.Instructions)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Items: results})
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}
