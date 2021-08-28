package get

import (
	"assessment/pkg/object"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RemoteObjectGetter struct {
	httpClient http.Client
	url        string
}

func NewRemoteObjectGetter(url string) *RemoteObjectGetter {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 200
	t.MaxConnsPerHost = 200
	t.MaxIdleConnsPerHost = 200
	t.IdleConnTimeout = 0

	return &RemoteObjectGetter{
		httpClient: http.Client{
			Timeout:   time.Second * 5,
			Transport: t,
		},
		url: url,
	}
}

func (r RemoteObjectGetter) Get(objectID int) (object.Object, error) {
	url := fmt.Sprintf("%s/objects/%d", r.url, objectID)

	resp, err := r.httpClient.Get(url)
	if err != nil {
		return object.Object{}, fmt.Errorf("error fetching object from url %s: %s", url, err)
	}
	defer resp.Body.Close()

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
