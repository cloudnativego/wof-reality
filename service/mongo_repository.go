package service

import (
	"errors"
	"fmt"

	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	"gopkg.in/mgo.v2/bson"
)

const (
	// RealityCollectionName defines the name of the MongoDB collection for storing game states
	RealityCollectionName = "realities"
)

type mongoReality struct {
	RecordID bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	GameID   string                 `json:"game_id"`
	GameMap  gameMap                `json:"game_map"`
	Players  map[string]playerState `json:"players"`
}

type mongoRealityRepository struct {
	Collection cfmgo.Collection
}

func newMongoRealityRepository(col cfmgo.Collection) (repo *mongoRealityRepository) {
	repo = &mongoRealityRepository{
		Collection: col,
	}
	return
}

func (repo *mongoRealityRepository) updateReality(gameID string, newReality reality) (err error) {
	repo.Collection.Wake()
	var recordID bson.ObjectId
	foundReality, err := repo.getRealityRecord(gameID)
	if err != nil {
		recordID = bson.NewObjectId()
	} else {
		recordID = foundReality.RecordID
	}
	fmt.Printf("Updating reality record, gameID: %s, recordID: %+v\n", newReality.GameID, recordID)
	newRecord := convertRealityToMongoReality(newReality, recordID)
	info, err := repo.Collection.UpsertID(recordID, newRecord)
	if err != nil {
		fmt.Printf("Failed to upsert record, %s\n", err.Error())
		return
	}
	if info.Updated != 1 {
		fmt.Printf("Failed to upsert: %+v\n", info)
		err = fmt.Errorf("Did not update 1 row, updated %d", info.Updated)
	}
	return
}

func (repo *mongoRealityRepository) getReality(gameID string) (gameReality reality, err error) {
	repo.Collection.Wake()
	record, err := repo.getRealityRecord(gameID)

	if err == nil {
		gameReality = convertMongoRealityToReality(record)
	}

	return
}

func (repo *mongoRealityRepository) getRealityRecord(gameID string) (realityRecord mongoReality, err error) {
	fmt.Printf("Looking up reality record for game %s\n", gameID)
	var records []mongoReality
	query := bson.M{"game_id": gameID}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := repo.Collection.Find(params, &records)
	if count == 0 {
		err = errors.New("Reality record not found.")
	}
	if err == nil {
		realityRecord = records[0]
	}
	return
}

func convertMongoRealityToReality(mongoRealityRecord mongoReality) (gameReality reality) {
	gameReality = reality{
		GameID:  mongoRealityRecord.GameID,
		GameMap: mongoRealityRecord.GameMap,
		Players: mongoRealityRecord.Players,
	}
	return
}

func convertRealityToMongoReality(gameReality reality, recordID bson.ObjectId) (mongoRealityRecord mongoReality) {
	mongoRealityRecord = mongoReality{
		RecordID: recordID,
		GameID:   gameReality.GameID,
		GameMap:  gameReality.GameMap,
		Players:  gameReality.Players,
	}
	return
}
