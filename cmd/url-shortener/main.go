package main 

import (
	"fmt"

	"github.com/finlleyl/shorty_reborn/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}