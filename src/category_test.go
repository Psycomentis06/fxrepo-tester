package src

import (
	"net/http"
	"testing"
)

func TestGetCategory(t *testing.T) {
  client := &http.Client{}
  GetCategory(client, "http://localhost:9057/api/v1/category/image/list", "")
}
