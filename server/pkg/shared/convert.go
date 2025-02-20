package shared

import "encoding/json"

func ConvertStruct(source any, destination any) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, destination); err != nil {
		return err
	}

	return nil
}
