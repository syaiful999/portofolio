package logs

import "encoding/json"

func ParseLogs(jsonData []byte) (*LogsResponse, error) {
	var response LogsResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
