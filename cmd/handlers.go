package functions

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	artist_groups string = "https://groupietrackers.herokuapp.com/api/artists"
	relation      string = "https://groupietrackers.herokuapp.com/api/relation"
)

type Error struct {
	CodeError        int
	ErrorDescription string
}

type filtersElement struct {
	CreationDateTo   string
	CreationDateFrom string
	FirstAlbumTo     string
	FirstAlbumFrom   string
	members          []string
	locations        string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	artists, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Artists, err := GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	tml := artists.Execute(w, Artists)
	if tml != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[:8] != "/artist/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(r.URL.Path[8:])
	if err != nil {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	if id > 52 || id <= 0 {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	artists, err := template.ParseFiles("./ui/html/artist-page.html")
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Artists, err := GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	tml := artists.Execute(w, Artists[id-1])
	if tml != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func Filters(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/filters/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	var res []Artist
	var value filtersElement
	Artists, err := GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	value.CreationDateTo = r.FormValue("CreationDateTo")
	value.CreationDateFrom = r.FormValue("CreationDateFrom")
	value.FirstAlbumTo = r.FormValue("FirstAlbumTo")
	value.FirstAlbumFrom = r.FormValue("FirstAlbumFrom")
	value.members = r.Form["member"]
	value.locations = r.FormValue("locations")
	AtoiCreationDateTo, _ := strconv.Atoi(value.CreationDateTo)
	AtoiCreationDateFrom, _ := strconv.Atoi(value.CreationDateFrom)
	AtoiFirstAlbumTo, _ := strconv.Atoi(value.FirstAlbumTo)
	AtoiFirstAlbumFrom, _ := strconv.Atoi(value.FirstAlbumFrom)
	membersInt := []int{}
	for _, v := range value.members {
		val, _ := strconv.Atoi(v)
		membersInt = append(membersInt, val)
	}
	if notNil(value.CreationDateTo, value.CreationDateFrom, value.FirstAlbumTo, value.FirstAlbumFrom, value.locations, value.members) {
		for _, v := range Artists {
			for _, member := range membersInt {
				for location := range v.DatesLocations {
					val, _ := strconv.Atoi(v.FirstAlbum[6:])
					for i := AtoiCreationDateTo; i <= AtoiCreationDateFrom; i++ {
						for j := AtoiFirstAlbumTo; j <= AtoiFirstAlbumFrom; j++ {
							if v.CreationDate == i && j == val && member == len(v.Members) && value.locations == location {
								res = append(res, v)
							} else {
								continue
							}
						}
					}
				}
			}
		}
	} else {
		for _, v := range Artists {
			if value.members == nil && value.locations != "" {
				for location := range v.DatesLocations {
					val, _ := strconv.Atoi(v.FirstAlbum[6:])
					for i := AtoiCreationDateTo; i <= AtoiCreationDateFrom; i++ {
						for j := AtoiFirstAlbumTo; j <= AtoiFirstAlbumFrom; j++ {
							if v.CreationDate == i && j == val && value.locations == location {
								res = append(res, v)
							} else {
								continue
							}
						}
					}
				}
			} else if value.members != nil && value.locations == "" {
				for _, member := range membersInt {
					val, _ := strconv.Atoi(v.FirstAlbum[6:])
					for i := AtoiCreationDateTo; i <= AtoiCreationDateFrom; i++ {
						for j := AtoiFirstAlbumTo; j <= AtoiFirstAlbumFrom; j++ {
							if v.CreationDate == i && j == val && member == len(v.Members) {
								res = append(res, v)
							} else {
								continue
							}
						}
					}
				}
			} else if value.members == nil && value.locations == "" {
				val, _ := strconv.Atoi(v.FirstAlbum[6:])
				for i := AtoiCreationDateTo; i <= AtoiCreationDateFrom; i++ {
					for j := AtoiFirstAlbumTo; j <= AtoiFirstAlbumFrom; j++ {
						if v.CreationDate == i && j == val {
							res = append(res, v)
						} else {
							continue
						}
					}
				}
			}
		}
	}
	tmpl, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result := tmpl.Execute(w, res)
	if result != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	var res []Artist
	query := strings.TrimSpace(r.FormValue("query"))
	if !checktxt(query) && query != "" {
		res = []Artist{}
	}
	Artists, err := GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	for _, v := range Artists {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(query)) && unique(res, v.Id) {
			res = append(res, v)
			continue
		}
		for _, members := range v.Members {
			if strings.Contains(strings.ToLower(members), strings.ToLower(query)) && unique(res, v.Id) {
				res = append(res, v)
				continue
			}
		}
		if strings.Contains(strconv.Itoa(v.CreationDate), query) && unique(res, v.Id) {
			res = append(res, v)
			continue
		}
		if strings.Contains(v.FirstAlbum, query) && unique(res, v.Id) {
			res = append(res, v)
			continue
		}
		for location := range v.DatesLocations {
			if strings.Contains(strings.ToLower(location), strings.ToLower(query)) && unique(res, v.Id) {
				res = append(res, v)
				continue
			}
		}
	}
	tmpl, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result := tmpl.Execute(w, res)
	if result != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func Errors(w http.ResponseWriter, errorNum int, errorDescript string) {
	tmpl, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errorNum)
	Error := Error{CodeError: errorNum, ErrorDescription: errorDescript}
	errors := tmpl.Execute(w, Error)
	if errors != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func checktxt(s string) bool {
	for _, v := range s {
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') {
			return true
		}
	}
	return false
}

func unique(res []Artist, i int) bool {
	for _, num := range res {
		if num.Id == i {
			return false
		}
	}
	return true
}

func notNil(CreationDateTo, CreationDateFrom, FirstAlbumTo, FirstAlbumFrom, location string, members []string) bool {
	for _, numOfMember := range members {
		if CreationDateTo != "" && FirstAlbumTo != "" && CreationDateFrom != "" && FirstAlbumFrom != "" && location != "" && numOfMember != "" {
			return true
		}
	}
	return false
}
