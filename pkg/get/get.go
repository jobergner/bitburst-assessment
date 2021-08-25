package get

import (
	"assessment/pkg/object"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func NewRemoteObjectGetter() RemoteObjectGetter {
	return RemoteObjectGetter{
		httpClient: http.Client{
			Timeout: time.Second * 5,
		},
	}
}

type RemoteObjectGetter struct {
	httpClient http.Client
}

func (r RemoteObjectGetter) Get(objectID int) (object.Object, error) {
	url := fmt.Sprintf("http://localhost:9091/objects/%d", objectID)

	resp, err := r.httpClient.Get(url)
	if err != nil {
		return object.Object{}, fmt.Errorf("error fetching object from url %s: %s", url, err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return object.Object{}, fmt.Errorf("error reading body of fetch response (object id: %d)", objectID)
	}

	var o object.Object
	err = json.Unmarshal(b, &o)
	if err != nil {
		return object.Object{}, fmt.Errorf("error unmarshalling object `%s`: %s", string(b), err)
	}

	return o, nil
}
