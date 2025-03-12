package api_test

//
//import (
//	"bytes"
//	"database/sql"
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"matcha/api/handlers"
//	"matcha/api/models"
//	"matcha/database"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//
//	_ "github.com/mattn/go-sqlite3"
//)
//
//// Configuration globale pour les tests
//var (
//	testDBPath  = "/tmp/matcha_chat_test.db"
//	testDB      *sql.DB
//	origInitDB  func() (*sql.DB, error)
//	testUserIDs []int
//)
//
//// TestMain configure l'environnement de test
//func TestMain(m *testing.M) {
//	// Sauvegarder la fonction InitDB originale
//	origInitDB = database.InitDB
//
//	// Configurer la BD avant les tests
//	if err := setupTestDB(); err != nil {
//		fmt.Printf("Erreur lors de la configuration de la BD de test: %v\n", err)
//		os.Exit(1)
//	}
//
//	// Exécuter les tests
//	code := m.Run()
//
//	// Nettoyer après tous les tests
//	cleanupTestDB()
//
//	os.Exit(code)
//}
//
//// setupTestDB initialise la base de données de test
//func setupTestDB() error {
//	// Supprimer la BD si elle existe déjà
//	os.Remove(testDBPath)
//
//	var err error
//	testDB, err = sql.Open("sqlite3", testDBPath)
//	if err != nil {
//		return fmt.Errorf("erreur lors de l'ouverture de la BD: %v", err)
//	}
//
//	// Remplacer la fonction InitDB pour qu'elle retourne notre BD de test
//	database.InitDB = func() (*sql.DB, error) {
//		return testDB, nil
//	}
//
//	// Charger le schéma
//	schema, err := ioutil.ReadFile("../database/schema.sql")
//	if err != nil {
//		return fmt.Errorf("impossible de lire le fichier schema.sql: %v", err)
//	}
//
//	// Créer les tables
//	if _, err = testDB.Exec(string(schema)); err != nil {
//		return fmt.Errorf("erreur lors de la création des tables: %v", err)
//	}
//
//	// Créer des utilisateurs de test
//	if err = createTestUsers(); err != nil {
//		return fmt.Errorf("erreur lors de la création des utilisateurs de test: %v", err)
//	}
//
//	return nil
//}
//
//// createTestUsers crée des utilisateurs pour les tests
//func createTestUsers() error {
//	// Créer 3 utilisateurs de test
//	users := []struct {
//		username string
//		email    string
//		password string
//		nom      string
//		prenom   string
//	}{
//		{"testuser1", "user1@test.com", "password1", "User1", "Test"},
//		{"testuser2", "user2@test.com", "password2", "User2", "Test"},
//		{"testuser3", "user3@test.com", "password3", "User3", "Test"},
//	}
//
//	testUserIDs = make([]int, len(users))
//
//	for i, user := range users {
//		result, err := testDB.Exec(`
//            INSERT INTO users (username, email, password, nom, prenom, first_step, terms_accepted)
//            VALUES (?, ?, ?, ?, ?, 1, 1)
//        `, user.username, user.email, user.password, user.nom, user.prenom)
//		if err != nil {
//			return err
//		}
//
//		id, err := result.LastInsertId()
//		if err != nil {
//			return err
//		}
//		testUserIDs[i] = int(id)
//	}
//
//	return nil
//}
//
//// resetTestDBData nettoie les données entre les tests sans recréer le schéma
//func resetTestDBData() error {
//	tables := []string{"messages", "chat_participants", "chat_rooms"}
//
//	for _, table := range tables {
//		if _, err := testDB.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
//			return err
//		}
//	}
//
//	// Réinitialiser les séquences d'autoincrement dans SQLite
//	if _, err := testDB.Exec("DELETE FROM sqlite_sequence WHERE name IN ('messages', 'chat_participants', 'chat_rooms')"); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// cleanupTestDB ferme la BD de test et supprime le fichier
//func cleanupTestDB() {
//	if testDB != nil {
//		testDB.Close()
//	}
//	database.InitDB = origInitDB
//	os.Remove(testDBPath)
//}
//
//// executeRequest est un helper pour exécuter des requêtes sur les handlers
//func executeRequest(req *http.Request, handler http.HandlerFunc) *httptest.ResponseRecorder {
//	rr := httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//	return rr
//}
//
//// Test pour le handler CreateChatroomHandler - cas normal
//func TestCreateChatroom(t *testing.T) {
//	// Réinitialiser les données de test
//	if err := resetTestDBData(); err != nil {
//		t.Fatalf("Erreur lors de la réinitialisation des données: %v", err)
//	}
//
//	// Préparer la requête
//	reqBody, err := json.Marshal(models.CreateChatroomRequest{
//		UserIDs: []int{testUserIDs[0], testUserIDs[1]},
//		Name:    "Test Chatroom",
//	})
//	if err != nil {
//		t.Fatalf("Erreur lors de la création du JSON de requête: %v", err)
//	}
//
//	req, err := http.NewRequest("POST", "/create_chatroom", bytes.NewBuffer(reqBody))
//	if err != nil {
//		t.Fatal(err)
//	}
//	req.Header.Set("Content-Type", "application/json")
//
//	// Exécuter la requête
//	rr := executeRequest(req, handlers.CreateChatroomHandler)
//
//	// Vérifier le code de statut
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("Handler a retourné un code status incorrect: got %v want %v, body: %s",
//			status, http.StatusOK, rr.Body.String())
//		return
//	}
//
//	// Vérifier la structure de la réponse
//	var response models.CreateChatroomResponse
//	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
//		t.Errorf("Impossible de décoder la réponse: %v", err)
//		return
//	}
//
//	// Vérifier le contenu de la réponse
//	if response.Status != "success" {
//		t.Errorf("Status incorrect dans la réponse: got %v want %v", response.Status, "success")
//	}
//	if response.ChatroomID <= 0 {
//		t.Errorf("ChatroomID invalide: got %v", response.ChatroomID)
//	}
//
//	// Vérifier que la chatroom a été créée dans la BD
//	var chatCount int
//	err = testDB.QueryRow("SELECT COUNT(*) FROM chat_rooms WHERE chatroom_id = ?", response.ChatroomID).Scan(&chatCount)
//	if err != nil {
//		t.Errorf("Erreur lors de la vérification de la chatroom: %v", err)
//	}
//	if chatCount != 1 {
//		t.Errorf("La chatroom n'a pas été créée en base de données")
//	}
//
//	// Vérifier que les participants ont été ajoutés
//	var partCount int
//	err = testDB.QueryRow("SELECT COUNT(*) FROM chat_participants WHERE chatroom_id = ?", response.ChatroomID).Scan(&partCount)
//	if err != nil {
//		t.Errorf("Erreur lors de la vérification des participants: %v", err)
//	}
//	if partCount != 2 {
//		t.Errorf("Les participants n'ont pas été correctement ajoutés. Attendu: 2, obtenu: %d", partCount)
//	}
//}
//
//// Test pour le handler CreateChatroomHandler - cas d'erreur: aucun utilisateur
//func TestCreateChatroomNoUsers(t *testing.T) {
//	if err := resetTestDBData(); err != nil {
//		t.Fatalf("Erreur lors de la réinitialisation des données: %v", err)
//	}
//
//	reqBody, _ := json.Marshal(models.CreateChatroomRequest{
//		UserIDs: []int{},
//		Name:    "Test Chatroom",
//	})
//
//	req, _ := http.NewRequest("POST", "/create_chatroom", bytes.NewBuffer(reqBody))
//	req.Header.Set("Content-Type", "application/json")
//
//	rr := executeRequest(req, handlers.CreateChatroomHandler)
//
//	// On attend un code 400 (Bad Request)
//	if status := rr.Code; status != http.StatusBadRequest {
//		t.Errorf("Handler a retourné un code status incorrect: got %v want %v",
//			status, http.StatusBadRequest)
//	}
//}
//
//// Test pour le handler CreateChatroomHandler - cas d'erreur: un seul utilisateur
//func TestCreateChatroomSingleUser(t *testing.T) {
//	if err := resetTestDBData(); err != nil {
//		t.Fatalf("Erreur lors de la réinitialisation des données: %v", err)
//	}
//
//	reqBody, _ := json.Marshal(models.CreateChatroomRequest{
//		UserIDs: []int{testUserIDs[0]},
//		Name:    "Test Chatroom",
//	})
//
//	req, _ := http.NewRequest("POST", "/create_chatroom", bytes.NewBuffer(reqBody))
//	req.Header.Set("Content-Type", "application/json")
//
//	rr := executeRequest(req, handlers.CreateChatroomHandler)
//
//	// On attend un code 400 (Bad Request)
//	if status := rr.Code; status != http.StatusBadRequest {
//		t.Errorf("Handler a retourné un code status incorrect: got %v want %v",
//			status, http.StatusBadRequest)
//	}
//}
//
//// Test pour le handler CreateChatroomHandler - cas d'erreur: utilisateur inexistant
//func TestCreateChatroomNonExistentUser(t *testing.T) {
//	if err := resetTestDBData(); err != nil {
//		t.Fatalf("Erreur lors de la réinitialisation des données: %v", err)
//	}
//
//	reqBody, _ := json.Marshal(models.CreateChatroomRequest{
//		UserIDs: []int{testUserIDs[0], 9999}, // ID invalide
//		Name:    "Test Chatroom",
//	})
//
//	req, _ := http.NewRequest("POST", "/create_chatroom", bytes.NewBuffer(reqBody))
//	req.Header.Set("Content-Type", "application/json")
//
//	rr := executeRequest(req, handlers.CreateChatroomHandler)
//
//	// On attend un code 404 (Not Found)
//	if status := rr.Code; status != http.StatusNotFound {
//		t.Errorf("Handler a retourné un code status incorrect: got %v want %v",
//			status, http.StatusNotFound)
//	}
//}
//
//// Test pour le handler CreateChatroomHandler - cas d'erreur: méthode non autorisée
//func TestCreateChatroomWrongMethod(t *testing.T) {
//	if err := resetTestDBData(); err != nil {
//		t.Fatalf("Erreur lors de la réinitialisation des données: %v", err)
//	}
//
//	req, _ := http.NewRequest("GET", "/create_chatroom", nil)
//	rr := executeRequest(req, handlers.CreateChatroomHandler)
//
//	// On attend un code 405 (Method Not Allowed)
//	if status := rr.Code; status != http.StatusMethodNotAllowed {
//		t.Errorf("Handler a retourné un code status incorrect: got %v want %v",
//			status, http.StatusMethodNotAllowed)
//	}
//}
