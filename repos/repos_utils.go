package repos

import (
	"errors"
	"github.com/softdev9/trendee-api-master/Godeps/_workspace/src/gopkg.in/mgo.v2/bson"
)

func BuildIdFromString(s string) (bson.ObjectId, error) {
	if bson.IsObjectIdHex(s) {
		id := bson.ObjectIdHex(s)
		return id, nil
	}
	return bson.NewObjectId(), errors.New("Invalid objects id given")
}

func IsValidBsonObjectId(idToCheck string) bool {
	return bson.IsObjectIdHex(idToCheck)
}
