package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"shorter/internal/config"
	"shorter/internal/dto"
	"shorter/internal/model"
	"shorter/internal/repository"
	"shorter/internal/utils"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ShoterHandler struct {
	repo   *repository.LinkRepository
	logger *zap.Logger
	cfg    *config.Config
}

func NewShorterHandler(
	repo *repository.LinkRepository,
	logger *zap.Logger,
	cfg *config.Config,
) *ShoterHandler {
	return &ShoterHandler{repo: repo, logger: logger, cfg: cfg}
}

type ErrResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	ErrMsg     string `json:"error"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.StatusCode)
	return nil
}

func ErrRender(err error, statusCode int) render.Renderer {
	return &ErrResponse{
		Err:        err,
		StatusCode: statusCode,
		ErrMsg:     err.Error(),
	}
}

func (s *ShoterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// request field to dto
	var req dto.ShorterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Render(w, r, ErrRender(err, http.StatusBadRequest))
		return
	}

	// validation
	var validate = validator.New()
	validate.RegisterValidation("url", func(fl validator.FieldLevel) bool {
		_, err := url.ParseRequestURI(fl.Field().String())
		return err == nil
	})
	if err := validate.Struct(req); err != nil {
		render.Render(w, r, ErrRender(err, http.StatusUnprocessableEntity))
		/*errs := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Tag()] = e.Error()
		}*/
		return
	}

	// generate alias
	alias := req.CustomAlias
	if alias == "" {
		randomAlias, err := utils.GenerateRandomAlias(6)
		if err != nil {
			s.logger.Fatal("failed to generate alias", zap.Error(err))
			render.Render(w, r, ErrRender(errors.New("internal error"), http.StatusInternalServerError))
			return
		}
		alias = randomAlias
	}

	// check alias already exists
	linkExist, err := s.repo.GetByAlias(r.Context(), alias)
	if err != nil {
		s.logger.Fatal("failed to get link by alias", zap.Error(err))
		render.Render(w, r, ErrRender(err, http.StatusInternalServerError))
		return
	}
	if linkExist != nil {
		render.Render(w, r, ErrRender(errors.New("alias already in use"), http.StatusConflict))
		return
	}

	// prepare model
	link := &model.Link{
		Alias:       alias,
		OriginalUrl: req.OriginalUrl,
	}

	if req.ExpiresIn != nil {
		duration := time.Duration(*req.ExpiresIn) * time.Second
		exp := time.Now().Add(duration)
		link.ExpiresAt = &exp
	}

	// save link
	if err := s.repo.Create(r.Context(), link); err != nil {
		s.logger.Error("failed to create link", zap.Error(err))
		render.Render(w, r, ErrRender(errors.New("internal error"), http.StatusInternalServerError))
		return
	}

	// response
	baseUrl := s.cfg.Server.Host + ":" + fmt.Sprint(s.cfg.Server.Port)
	resp := dto.ShorterResponse{
		ShortUrl: baseUrl + "/" + alias,
	}
	if link.ExpiresAt != nil {
		resp.ExpiresAt = link.ExpiresAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
