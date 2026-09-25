package callouts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var getCache *cache = NewCache(10 * time.Second)

func DoGetCall[T any](url string, parsedResponse *T) error {
	resBody, found := getCache.Get(url)

	if !found {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return fmt.Errorf("Failed to Get with  StatusCode %v, Status %v.", res.StatusCode, res.Status)
		}

		resBody, err = io.ReadAll(res.Body)
		getCache.Add(url, resBody)
		if err != nil {
			return err
		}
	}

	if err := json.Unmarshal(resBody, &parsedResponse); err != nil {
		return err
	}

	return nil
}
