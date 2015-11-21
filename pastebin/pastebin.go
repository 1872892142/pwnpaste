/* pastebin is a wrapper for pastebin.com. */

package pastebin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetPaste(pasteID string) ([]byte, error) {
	var data []byte

	url := fmt.Sprintf("https://pastebin.com/raw.php?i=%s", pasteID)
	res, err := http.Get(url)
	if err != nil {
		return data, err
	}
	if res.StatusCode == http.StatusNotFound {
		return data, errors.New("404 Not Found")
	}

	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	return data, err
}
