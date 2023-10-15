package src

import (
	"net/http"
	"testing"
)

func TestGetCategory(t *testing.T) {
	client := &http.Client{}
	serverHost := "http://localhost:9057"
	getImageEndpoint := serverHost + "/api/v1/category/"
	categoryNames := []string{"default", "cars", "landscape", "cats"}
	for _, categoryName := range categoryNames {
		cat, catErr := GetCategory(client, getImageEndpoint+categoryName)
		if catErr != nil {
			t.Errorf("Failed to Get Category \"%s\": %s", categoryName, catErr)
		} else {
			t.Log("Category Found: ", cat.Name)
		}
	}
}

func TestCreateCategory(t *testing.T) {
	client := &http.Client{}
	serverHost := "http://localhost:9057"
	createImageEndpoint := serverHost + "/api/v1/category/new"
	category, err := CreateCategory(client, createImageEndpoint, Category{Name: "Test#1", Description: "Test Description"})
	if err != nil {
		t.Errorf("Error creating category: %s ", err)
	} else {
		t.Log("Category Created: ", category.Name)
	}
}
