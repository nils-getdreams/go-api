package api

import (
	"fmt"
	"strconv"
	"time"
)

func (api *API) waitUntilReady(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	for {
		resp, err := api.sling.Path("/api/instances/").Get(id).ReceiveSuccess(&data)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Got statuscode %d from api: %s", resp.StatusCode, resp.Status)
		}
		if data["ready"] == true {
			data["id"] = id
			return data, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (api *API) CreateInstance(params map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	resp, err := api.sling.Post("/api/instances").BodyJSON(params).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got statuscode %d from api: %s", resp.StatusCode, resp.Status)
	}

	if _, ok := data["id"]; !ok {
		return nil, fmt.Errorf("No id in response")
	}
	string_id := strconv.Itoa(int(data["id"].(float64)))
	return api.waitUntilReady(string_id)
}

func (api *API) ReadInstance(id string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	resp, err := api.sling.Path("/api/instances/").Get(id).ReceiveSuccess(&data)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Got statuscode %d from api: %s ", resp.StatusCode, resp.Status)
	}
	return data, nil
}

func (api *API) UpdateInstance(id string, params map[string]interface{}) error {
	resp, err := api.sling.Put("/api/instances/" + id).BodyJSON(params).ReceiveSuccess(nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Got statuscode %d from api: %s", resp.StatusCode, resp.Status)
	}
	return err
}

func (api *API) DeleteInstance(id string) error {
	resp, err := api.sling.Path("/api/instances/").Delete(id).ReceiveSuccess(nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("Got statuscode %d from api: %s", resp.StatusCode, resp.Status)
	}
	return err
}
