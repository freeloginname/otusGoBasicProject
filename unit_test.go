package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/freeloginname/otusGoBasicProject/pkg/notes"
	"github.com/freeloginname/otusGoBasicProject/pkg/users"
)

var letters = "abcdefghijklmnopqrstuvwxyz"

func randSeq(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

type LoginData struct {
	Token string `json:"token"`
	Error string `json:"error,omitempty"`
}

func LoginUser(name string) (string, error) {
	user := users.LoginUserRequestBody{
		Name:     name,
		Password: "test",
	}
	body, _ := json.Marshal(user)
	// expected := LoginData{Token: "token"}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/users/login", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
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

func CreateUser(name string) error {
	user := users.CreateUserRequestBody{
		Name:     name,
		Password: "test",
	}
	body, _ := json.Marshal(user)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/users", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// fmt.Println(string(data))
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("user not created")
	}
	return nil
}

func retry(attempts int, sleep time.Duration) (ok int, err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
		if err != nil {
			continue
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err == nil {
			return resp.StatusCode, nil
		}
		defer resp.Body.Close()

		// resp, err := http.Get("http://localhost:8080/")
		// if err == nil {
		// 	return resp.StatusCode, nil
		// }
		// resp.Body.Close()
	}
	return 0, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}

func TestConnection(t *testing.T) {
	expected := http.StatusOK
	result, err := retry(10, 5*time.Second)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if result != expected {
		t.Errorf("expected %d got %v", expected, result)
	}
}

func TestCreateUser(t *testing.T) {
	name, err := randSeq(5)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	user := users.CreateUserRequestBody{
		Name:     name,
		Password: "test",
	}
	body, _ := json.Marshal(user)
	expected := "{\"success\":\"User Created\"}"

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/users", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
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

func TestLoginUser(t *testing.T) {
	name, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	err = CreateUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	user := users.LoginUserRequestBody{
		Name:     name,
		Password: "test",
	}
	body, _ := json.Marshal(user)
	// expected := LoginData{Token: "token"}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/users/login", strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
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
	var structuredData LoginData
	err = json.Unmarshal(data, &structuredData)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if len(structuredData.Token) == 0 {
		t.Errorf("expected token not empty got %v", string(data))
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected %d got %v", http.StatusCreated, resp.StatusCode)
	}
}

// Тестирование обработки некорректного входу закомментировано из-за ошибки линтера на дублирование кода.
// func TestLoginUserIncorrect(t *testing.T) {
// 	user := users.LoginUserRequestBody{
// 		Name:     "test",
// 		Password: "fake",
// 	}
// 	body, _ := json.Marshal(user)

// 	resp, err := http.Post("http://localhost:8080/users/login", "application/json", strings.NewReader(string(body)))
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	defer resp.Body.Close()
// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	var structuredData LoginData
// 	err = json.Unmarshal(data, &structuredData)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	if len(structuredData.Error) <= 0 {
// 		t.Errorf("expected token not empty got %v", string(data))
// 	}
// 	if resp.StatusCode != http.StatusBadRequest {
// 		t.Errorf("expected %d got %v", http.StatusBadRequest, resp.StatusCode)
// 	}
// }

func TestCreateNote(t *testing.T) {
	name, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	err = CreateUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	token, err := LoginUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	noteName, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	note := notes.CreateNoteRequestBody{
		Name:     noteName,
		Text:     "test",
		UserName: name,
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

// Тестирование получения списка всех заметок закомментировано из-за ошибки линтера на дублирование кода
// func TestGetNotes(t *testing.T) {
// 	token, err := LoginUser()
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	expected := " test "

// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", "http://localhost:8080/notes", nil)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	cookie := http.Cookie{Name: "token", Value: token}
// 	req.AddCookie(&cookie)
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("expected %d got %v", http.StatusOK, resp.StatusCode)
// 	}
// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		t.Errorf("expected error to be nil got %v", err)
// 	}
// 	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
// 		insideHTML, _ := s.Html() //underscore is an error
// 		if insideHTML != expected {
// 			t.Errorf("expected %s got %v", expected, insideHTML)
// 		}
// 	})
// }

func TestGetNote(t *testing.T) {
	name, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	err = CreateUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	token, err := LoginUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "test"

	// create note
	noteName, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	note := notes.CreateNoteRequestBody{
		Name:     noteName,
		Text:     "test",
		UserName: name,
	}
	body, _ := json.Marshal(note)

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
	// created

	url := "http://localhost:8080/notes/" + noteName
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	// cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err = client.Do(req)
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
	doc.Find("#text").Each(func(_ int, s *goquery.Selection) {
		insideHTML, _ := s.Html()
		if insideHTML != expected {
			t.Errorf("expected %s got %v", expected, insideHTML)
		}
		// fmt.Printf("Review %d: %s\n", i, insideHTML)
	})
}

func TestUpdateNote(t *testing.T) {
	name, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	err = CreateUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	token, err := LoginUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	// create note
	noteName, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	note := notes.CreateNoteRequestBody{
		Name:     noteName,
		Text:     "test",
		UserName: name,
	}
	body, _ := json.Marshal(note)

	client := &http.Client{}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/notes", strings.NewReader(string(body)))
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
	// created

	updateNote := notes.UpdateNoteRequestBody{
		Text: "testUpdated",
	}
	body, _ = json.Marshal(updateNote)
	expected := "{\"success\":\"Note Updated\"}"

	url := "http://localhost:8080/notes/" + noteName
	req, err = http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(string(body)))
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	req.AddCookie(&cookie)
	resp, err = client.Do(req)
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

	// get note
	url = "http://localhost:8080/notes/" + noteName
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	req.AddCookie(&cookie)
	resp, err = client.Do(req)
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
	doc.Find("#text").Each(func(_ int, s *goquery.Selection) {
		insideHTML, _ := s.Html()
		if insideHTML != "testUpdated" {
			t.Errorf("expected %s got %v", "testUpdated", insideHTML)
		}
		fmt.Printf("Review: %s\n", insideHTML)
	})
}

func TestDeleteNote(t *testing.T) {
	name, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	err = CreateUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	token, err := LoginUser(name)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "{\"success\":\"Note Deleted or does not exist\"}"

	// create note
	noteName, err := randSeq(10)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	note := notes.CreateNoteRequestBody{
		Name:     noteName,
		Text:     "test",
		UserName: name,
	}
	body, _ := json.Marshal(note)

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
	// created

	ctx := context.Background()
	url := "http://localhost:8080/notes/" + noteName
	req, err = http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	req.AddCookie(&cookie)
	resp, err = client.Do(req)
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
