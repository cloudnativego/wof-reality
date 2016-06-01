package integrations_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/cloudnativego/wof-reality/service"
)

var (
	appEnv, _ = cfenv.Current()
	server    = NewServer(appEnv)
)

func TestIntegration(t *testing.T) {
	fmt.Println("Integration test placeholder")

	_, err := getReality("noneexistent", http.StatusNotFound, t)
	if err != nil {
		t.Errorf("Should have gotten 404 but didn't: %s", err.Error())
	}

	// TODO: make integration test pass for upserting.
	newReality := createFakeReality("splendid")
	err = putReality("splendid", newReality, t)
	if err != nil {
		t.Errorf("Failed to create reality: %s", err.Error())
		return
	}

	expectedGame, err := getReality("splendid", http.StatusOK, t)
	if err != nil {
		t.Errorf("Should've gotten 200 when retrieving existing game, got %s", err.Error())
	}
	if expectedGame.Players["bob"].Hitpoints != 99 {
		t.Errorf("Retrieved game doesn't match what we expected, %+v", expectedGame)
	}

	// change a game state and update it.
	newReality.Players["bob"] = playerState{Hitpoints: 1, ID: "bob", CurrentTileID: "tile1"}
	newReality.GameMap.Metadata.Author = "Changey McChangePants"
	err = putReality("splendid", newReality, t)
	if err != nil {
		t.Errorf("Failed to update reality: %s", err.Error())
		return
	}

	updatedGame, err := getReality("splendid", http.StatusOK, t)
	if err != nil {
		t.Errorf("Should've gotten 200 when retrieving updated game, got %s", err.Error())
		return
	}

	if updatedGame.Players["bob"].Hitpoints != 1 {
		t.Errorf("Bob's hitpoints should have been updated to 1, got %d", updatedGame.Players["bob"].Hitpoints)
	}

	if updatedGame.GameMap.Metadata.Author != "Changey McChangePants" {
		t.Errorf("Map author should have changed but didn't, author was %s", updatedGame.GameMap.Metadata.Author)
	}
}

/* ========================== */

func getReality(gameID string, expectedCode int, t *testing.T) (gameReality reality, err error) {
	getRealityRequest, _ := http.NewRequest("GET", fmt.Sprintf("/reality/%s", gameID), nil)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, getRealityRequest)

	if expectedCode != 200 && expectedCode != 201 { // skip the de-serialize if we don't expect a variant of OK
		return
	}

	err = json.Unmarshal(recorder.Body.Bytes(), &gameReality)
	if err != nil {
		t.Errorf("Error unmarshaling game reality, %v", err)
	} else {
		if recorder.Code != expectedCode {
			t.Errorf("Expected reality query code to be %d, got %d", expectedCode, recorder.Code)
		} else {
			fmt.Println("\tQueried Reality OK")
		}
	}
	return
}

func putReality(gameID string, gameReality reality, t *testing.T) (err error) {
	rawbytes, _ := json.Marshal(gameReality)
	reader := bytes.NewBuffer(rawbytes)
	putRealityRequest, _ := http.NewRequest("PUT", fmt.Sprintf("/reality/%s", gameID), reader)
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, putRealityRequest)

	if recorder.Code != http.StatusOK {
		t.Errorf("Failed to PUT reality state, got code %d", recorder.Code)
		err = errors.New("Failed to PUT reality state")
	}
	return
}

// This may look like a lot of duplication, but, as we discuss in the book, it's
// better than creating tight coupling.

func createFakeReality(gameID string) (gameReality reality) {
	gameReality = reality{
		GameID:  gameID,
		GameMap: createFakeMap(),
		Players: createFakePlayerState(),
	}
	return
}

func createFakeMap() (fakeMap gameMap) {
	fakeMap = gameMap{
		Tiles: make([][]mapTile, 1),
		ID:    "testmap1",
		Metadata: mapMetadata{
			Author:      "Test Map Maker",
			Description: "Test Map",
		},
	}
	fakeMap.Tiles[0] = make([]mapTile, 1)
	fakeMap.Tiles[0][0] = createFakeMapTile()
	return
}

func createFakeMapTile() (tile mapTile) {
	tile = mapTile{
		ID:      "tile1",
		Sprite:  "",
		AllowUp: true, AllowDown: true, AllowLeft: true, AllowRight: true,
		Traversable: true, TileName: "grass-dirt-12",
	}
	return
}

func createFakePlayerState() (players map[string]playerState) {
	players = make(map[string]playerState)
	players["bob"] = playerState{
		ID:            "bob",
		CurrentTileID: "tile1",
		Hitpoints:     99,
		Name:          "Bob",
		Sprite:        "elf-1",
	}
	return
}

type gameMap struct {
	Tiles    [][]mapTile `json:"tiles"`
	ID       string      `json:"id"`
	Metadata mapMetadata `json:"metadata"`
}

type mapMetadata struct {
	Author      string `json:"author"`
	Description string `json:"description"`
}

type mapTile struct {
	ID          string `json:"id"`
	Sprite      string `json:"sprite"`
	AllowUp     bool   `json:"allow_up"`
	AllowDown   bool   `json:"allow_down"`
	AllowLeft   bool   `json:"allow_left"`
	AllowRight  bool   `json:"allow_right"`
	Traversable bool   `json:"traversable"`
	TileName    string `json:"tile_name"`
}

type playerState struct {
	Hitpoints     uint   `json:"hit_points"`
	ID            string `json:"player_id"`
	CurrentTileID string `json:"current_tile_id"`
	Name          string `json:"name"`
	Sprite        string `json:"sprite"`
}

type reality struct {
	GameID  string                 `json:"game_id"`
	GameMap gameMap                `json:"game_map"`
	Players map[string]playerState `json:"players"`
}
