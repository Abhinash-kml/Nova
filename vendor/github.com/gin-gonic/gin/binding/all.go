package binding

import (
	"net/http"

	"github.com/gin-gonic/gin/codec/json"
)

type allSourceBinding struct{}

func (ab allSourceBinding) Name() string {
	return "all"
}

func (ab allSourceBinding) BindAll(req *http.Request, m map[string][]string, obj any) error {
	// Header param binding
	if len(req.Header) > 0 {
		if err := mapHeader(obj, req.Header); err != nil {
			return err
		}
	}

	// Uri param binding
	if req.URL != nil {
		if err := mapURI(obj, m); err != nil {
			return err
		}
	}

	// Query and form param binding
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}

	// Body
	// b := Default(req.Method, req.Header.Get("Content-Type"))
	if req.Body != nil && req.Header.Get("Content-Type") == "application/json" {
		decoder := json.API.NewDecoder(req.Body)
		if EnableDecoderUseNumber {
			decoder.UseNumber()
		}
		if EnableDecoderDisallowUnknownFields {
			decoder.DisallowUnknownFields()
		}
		if err := decoder.Decode(obj); err != nil {
			return err
		}
	}

	return validate(obj)
}
