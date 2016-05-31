package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateReality(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	fakeRepo := newInMemoryRepository()
	server := newServerWithRepo(fakeRepo)
	recorder = httptest.NewRecorder()

	fakeReality := createFakeReality("game90")
	realityBytes, _ := json.Marshal(fakeReality)

	reader := bytes.NewReader(realityBytes)
	request, _ = http.NewRequest("PUT", "/reality/game90", reader)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Updating game state, we should have received a HTTP 200, got %d", recorder.Code)
	}

	if _, ok := fakeRepo.states["game90"]; !ok {
		t.Errorf("No game state was stored in repository after HTTP PUT.")
	}

	state, _ := fakeRepo.states["game90"]
	if state.GameMap.Tiles[0][0].ID != "tile1" {
		t.Errorf("Lost information during update to the repository.")
	}
	if state.Players["bob"].Hitpoints != 99 {
		t.Errorf("Bob should have had 99 hitpoints, instead had %d", state.Players["bob"].Hitpoints)
	}

	fakeReality.Players["bob"] = playerState{ID: "bob", Hitpoints: 1}
	realityBytes2, _ := json.Marshal(fakeReality)
	reader = bytes.NewReader(realityBytes2)
	recorder = httptest.NewRecorder()
	request2, _ := http.NewRequest("PUT", "/reality/game90", reader)
	server.ServeHTTP(recorder, request2)

	if recorder.Code != http.StatusOK {
		t.Errorf("Should've gotten an OK on 2nd update, got %d", recorder.Code)
	}
	state2, _ := fakeRepo.states["game90"]
	if state2.Players["bob"].Hitpoints != 1 {
		t.Errorf("Should have reduced bob's hitpoints during an update to 1, instead got %d\n", state2.Players["bob"].Hitpoints)
	}

	return
}

func TestGetReality(t *testing.T) {
	fmt.Println("placeholder")
	return
}

// ====================

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
