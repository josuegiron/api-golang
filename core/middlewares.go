package core

import (
	"net/http"
	"net/http/httptest"
)

// Middleware provides a convenient mechanism for filtering HTTP requests
// entering the application. It returns a new handler which may perform various
// operations and should finish by calling the next HTTP handler.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// MiddlewaresChain provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func MiddlewaresChain(mw ...Middleware) Middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			res := httptest.NewRecorder()
			// printAPIRequest(r)
			last(res, r)
			for i, header := range res.Header() {
				w.Header()[i] = header
			}
			w.WriteHeader(res.Code)
			w.Write(res.Body.Bytes())
			// printAPIResponse(res)
		}
	}
}

// func printAPIRequest(r *http.Request) {
// 	byteBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Infof("REQUEST: %v %v", r.Method, r.RequestURI)
// 	log.Infof("HEADER: %v", r.Header)
// 	log.Infof("FORM: %v", r.Form.Encode())
// 	log.Infof("BODY: \n%v", string(byteBody))
// 	r.Body = ioutil.NopCloser(bytes.NewBuffer(byteBody))
// }

// func printAPIResponse(res *httptest.ResponseRecorder) {
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(res.Body)
// 	bodyStr := buf.String()
// 	log.Infof("RESPONSE: %v", res.Result())
// 	log.Infof("STATUS CODE: %v %v", res.Code, http.StatusText(res.Code))
// 	log.Infof("HEADER: %v", res.Header())
// 	log.Infof("BODY: \n%v", bodyStr)
// }

// // ValidateRequest doc ...
// func (api APIRest) ValidateRequest(r *http.Request) error {
// 	err := api.UnmarshalBody(&api.Request, r)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }