package vault

import "encoding/json"

func Marshal(v *Vault) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte) (*Vault, error) {
	var v Vault
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
