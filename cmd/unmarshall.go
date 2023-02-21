package functions

import (
	"encoding/json"
	"io"
	"net/http"
)

type Artist struct {
	Id             int
	Image          string
	Name           string
	Members        []string
	CreationDate   int
	FirstAlbum     string
	DatesLocations map[string][]string
}

type Relation struct {
	Index []struct {
		Id             uint64
		DatesLocations map[string][]string
	}
}

var (
	Art []Artist
	Rel Relation
)

func GetAllArtist(relationArtist, Artistlink string) ([]Artist, error) {
	res, err := http.Get(Artistlink)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, errBody := io.ReadAll(res.Body)
	if errBody != nil {
		return nil, errBody
	}
	if jsonErr := json.Unmarshal(body, &Art); jsonErr != nil {
		return nil, jsonErr
	}
	resRel, err := http.Get(relationArtist)
	if err != nil {
		return nil, err
	}
	defer resRel.Body.Close()
	bodyRel, errBody := io.ReadAll(resRel.Body)
	if errBody != nil {
		return nil, errBody
	}
	if jsonErrRel := json.Unmarshal(bodyRel, &Rel); jsonErrRel != nil {
		return nil, jsonErrRel
	}
	for i := range Art {
		Art[i].DatesLocations = Rel.Index[i].DatesLocations
	}
	return Art, nil
}
