package middlewares

import (
	"net/http"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

var ErrValidationFailed = errors.Template().With(
	opts.Title("Validation failed"),
	opts.StatusCode(http.StatusBadRequest),
)

func ValidatedInput(inputStruct any) func(http.Handler) http.Handler {
	inputMiddleware := httpin.NewInput(inputStruct)
	validate := validator.New()
	return func(next http.Handler) http.Handler {
		return inputMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			input := r.Context().Value(httpin.Input)
			if err := validate.Struct(input); err != nil {
				errs := err.(validator.ValidationErrors)
				rjson.HandleError(w, ErrValidationFailed.With(rjson.NewValidationError(errs)))
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

func init() {
	httpin.UseGochiURLParam("path", chi.URLParam)
}
