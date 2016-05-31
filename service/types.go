package service

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

type realityRepository interface {
	updateReality(gameID string, newReality reality) (err error)
	getReality(gameID string) (gameReality reality, err error)
}
