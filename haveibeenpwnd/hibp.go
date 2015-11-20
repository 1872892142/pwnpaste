// Author:  Jonathan Broche and Tom Steele

package hibp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type PasteAccount []struct {
	Source     string    `json:"Source"`
	ID         string    `json:"Id"`
	Title      string    `json:"Title"`
	Date       time.Time `json:"Date"`
	Emailcount int       `json:"EmailCount"`
}

func GetPasteAccount(email string) (PasteAccount, error) {
	var p PasteAccount
	url := fmt.Sprintf("https://haveibeenpwned.com/api/v2/pasteaccount/%s", email)
	res, err := http.Get(url)
	if err != nil {
		return p, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return p, err
	}
	res.Body.Close()

	err = json.Unmarshal(body, &p)
	return p, err
}
