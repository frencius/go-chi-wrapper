package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDoc(t *testing.T) {
	mockSwagger := Swagger{}

	t.Run("negative", func(t *testing.T) {
		expectedSwag := "{\n    \"info\": {\n        \"title\": \"Swagger Test\"\n    }\n}"

		mockFilePath := "swagger_test.json"
		Path = mockFilePath

		swagStr := mockSwagger.ReadDoc()
		assert.Equal(t, swagStr, expectedSwag)
	})

	t.Run("negative - read file error", func(t *testing.T) {
		expectedSwag := "open : no such file or directory"

		mockFilePath := ""
		Path = mockFilePath

		swagStr := mockSwagger.ReadDoc()
		assert.Equal(t, swagStr, expectedSwag)
	})
}
