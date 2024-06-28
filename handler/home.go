package handler

import (
	"genartai/view/home"
	"net/http"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	//return fmt.Errorf("failed to load picture")
	// user := getAuthenticatedUser(r)
	// account, err := db.GetAccountByUserID(user.ID)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("%+v\n", user.Account)
	return home.Index().Render(r.Context(), w)
}
