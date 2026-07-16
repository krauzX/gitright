package main

import (
	"fmt"
	"os"
	"time"

	"github.com/krauzx/gitright/internal/middleware"
)

func main() {
	secret := "/7xGdY2CkVJc80ELqnOlhxPvCrFats4xNKCBFoUHKvk="
	token, err := middleware.GenerateJWT(1, "krauzX", secret, 24*time.Hour)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(token)
}
