package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/freeloginname/otusGoBasicProject/pkg/notes"
	"github.com/freeloginname/otusGoBasicProject/pkg/users"

	"github.com/PuerkitoBio/goquery"
)

type LoginData struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func LoginUser() (string, error) {
	user := users.LoginUserRequestBody{
		Name:     "test",
		Password: "test",
	}
	body, _ := json.Marshal(user)
	// expected := LoginData{Token: "token"}

	resp, err := http.Post("http://localhost:8080/users/login", "application/json", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// fmt.Println(string(data))
	var structuredData LoginData
	err = json.Unmarshal(data, &structuredData)
	if err != nil {
		return "", err
	}
	return structuredData.Token, nil
}

func TestCreateUser(t *testing.T) {
	user := users.CreateUserRequestBody{
		Name:     "test",
		Password: "test",
	}
	body, _ := json.Marshal(user)
	expected := "{\"success\":\"User Created\"}"

	resp, err := http.Post("http://localhost:8080/users", "application/json", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// fmt.Println(string(data))
	if string(data) != expected {
		t.Errorf("expected %s got %v", expected, string(data))
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestLoginUser(t *testing.T) {
	user := users.LoginUserRequestBody{
		Name:     "test",
		Password: "test",
	}
	body, _ := json.Marshal(user)
	// expected := LoginData{Token: "token"}

	resp, err := http.Post("http://localhost:8080/users/login", "application/json", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// fmt.Println(string(data))
	var structuredData LoginData
	err = json.Unmarshal(data, &structuredData)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if len(structuredData.Token) <= 0 {
		t.Errorf("expected token not empty got %v", string(data))
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestLoginUserIncorrect(t *testing.T) {
	user := users.LoginUserRequestBody{
		Name:     "test",
		Password: "fake",
	}
	body, _ := json.Marshal(user)
	// expected := LoginData{Token: "token"}

	resp, err := http.Post("http://localhost:8080/users/login", "application/json", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// fmt.Println(string(data))
	var structuredData LoginData
	err = json.Unmarshal(data, &structuredData)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if len(structuredData.Error) <= 0 {
		t.Errorf("expected token not empty got %v", string(data))
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected %d got %v", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestCreateNote(t *testing.T) {
	token, err := LoginUser()
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	note := notes.CreateNoteRequestBody{
		Name:     "test",
		Text:     "test",
		UserName: "test",
	}
	body, _ := json.Marshal(note)
	expected := "{\"success\":\"Note Created\"}"

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/notes", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// fmt.Println(string(data))
	if string(data) != expected {
		t.Errorf("expected %s got %v", expected, string(data))
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d got %v", http.StatusCreated, resp.StatusCode)
	}
}

func TestGetNotes(t *testing.T) {
	token, err := LoginUser()
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := " test "

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/notes", nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	//
	//data, err := io.ReadAll(resp.Body)
	//fmt.Println(string(data))
	//
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %v", http.StatusOK, resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		inside_html, _ := s.Html() //underscore is an error
		if inside_html != expected {
			t.Errorf("expected %s got %v", expected, inside_html)
		}
		// fmt.Printf("Review %d: %s\n", i, inside_html)
	})
}

func TestGetNote(t *testing.T) {
	token, err := LoginUser()
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "test"

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/notes/test", nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %v", http.StatusOK, resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	doc.Find("#text").Each(func(i int, s *goquery.Selection) {
		inside_html, _ := s.Html() //underscore is an error
		if inside_html != expected {
			t.Errorf("expected %s got %v", expected, inside_html)
		}
		// fmt.Printf("Review %d: %s\n", i, inside_html)
	})
}

func TestUpdateNote(t *testing.T) {
	token, err := LoginUser()
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	note := notes.UpdateNoteRequestBody{
		Text: "testUpdated",
	}
	body, _ := json.Marshal(note)
	expected := "{\"success\":\"Note Updated\"}"

	client := &http.Client{}
	req, err := http.NewRequest("PUT", "http://localhost:8080/notes/test", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// fmt.Println(string(data))
	if string(data) != expected {
		t.Errorf("expected %s got %v", expected, string(data))
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d got %v", http.StatusCreated, resp.StatusCode)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/notes/test", nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	req.AddCookie(&cookie)
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %d got %v", http.StatusOK, resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	doc.Find("#text").Each(func(i int, s *goquery.Selection) {
		inside_html, _ := s.Html() //underscore is an error
		if inside_html != "testUpdated" {
			t.Errorf("expected %s got %v", "testUpdated", inside_html)
		}
		fmt.Printf("Review %d: %s\n", i, inside_html)
	})
}

func TestDeleteNote(t *testing.T) {
	token, err := LoginUser()
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "{\"success\":\"Note Deleted or does not exist\"}"

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://localhost:8080/notes/test", nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != expected {
		t.Errorf("expected %s got %v", expected, string(data))
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d got %v", http.StatusCreated, resp.StatusCode)
	}
}
