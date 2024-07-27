package parsing

import "encoding/json"

func MustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b
}
