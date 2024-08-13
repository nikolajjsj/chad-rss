package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateNID() (string, error) {
	return gonanoid.New()
}
