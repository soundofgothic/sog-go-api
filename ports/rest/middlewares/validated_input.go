package middlewares

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

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

type IDArray struct {
	Values []int64
}

func decodeIDArray(value string) (any, error) {
	if value == "" {
		return IDArray{}, nil
	}
	strIDs := strings.Split(value, ",")
	ids := make([]int64, len(strIDs))
	for i, strID := range strIDs {
		id, err := strconv.ParseInt(strID, 10, 64)
		if err != nil {
			return nil, errors.Enhance(err)
		}
		ids[i] = id
	}
	return IDArray{Values: ids}, nil
}

func init() {
	httpin.UseGochiURLParam("path", chi.URLParam)
	httpin.RegisterTypeDecoder(reflect.TypeOf(IDArray{}), httpin.ValueTypeDecoderFunc(decodeIDArray))
}
