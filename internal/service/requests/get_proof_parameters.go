package requests

//type UserInputs struct {
//	UserId        string `url:"userId"`
//	AgeLowerBound int    `url:"ageLowerBound"`
//	Uniqueness    bool   `url:"uniqueness"`
//}

//func NewGetUserInputs(r *http.Request) (userInputs UserInputs, err error) {
//if err = urlval.Decode(r.URL.Query(), &userInputs); err != nil {
//	err = newDecodeError("query", err)
//	return
//}
//return
//}
//
//func newDecodeError(what string, err error) error {
//	return validation.Errors{
//		what: fmt.Errorf("decode request %s: %w", what, err),
//	}
//}
