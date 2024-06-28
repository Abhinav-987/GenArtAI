package handler

import (
	"genartai/view/credits"
	"net/http"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, credits.Index())
}
