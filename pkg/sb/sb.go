package sb

import (
	"errors"
	"os"

	"github.com/nedpals/supabase-go"
)

var BaseAuthURL = "https://rzyghlyzrlyrendtzvec.supabase.co/auth/v1/recover"

var Client *supabase.Client

func Init() error {
	sbHost := os.Getenv("SUPABASE_URL")
	if sbHost == "" {
		return errors.New("supabase host is required")
	}
	sbSecret := os.Getenv("SUPABASE_SECRET")
	if sbSecret == "" {
		return errors.New("supabase secret is required")
	}
	Client = supabase.CreateClient(sbHost, sbSecret)
	return nil
}
