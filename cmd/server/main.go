package server

import "honnef.co/go/tools/config"

func main() {
	cfg := config.Load()

	db := postgres.New(cfg.DB)
	router := api.NewRouter(db, cfg)

	router.Run(":8080")
}
