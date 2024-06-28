package handler

import (
	"net/http"

	"github.com/Abhinav-987/GenArtAI/view/credits"
)

func HandleCreditsIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, credits.Index())
}
