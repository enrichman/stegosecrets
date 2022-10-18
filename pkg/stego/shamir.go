package stego

import (
	"encoding/base64"

	shamir "github.com/corvus-ch/shamir"
)

type Part struct {
	Version byte
	Tag     byte
	Content []byte
}

func NewPart(content []byte) Part {
	return Part{
		Version: content[0],
		Tag:     content[1],
		Content: content[2:],
	}
}

func (p Part) Bytes() []byte {
	return append([]byte{
		p.Version,
		p.Tag,
	}, p.Content...)
}

func (p Part) Base64() string {
	return base64.StdEncoding.EncodeToString(p.Bytes())
}

func Split(secret []byte, parts, threshold int) ([]Part, error) {
	partsMap, err := shamir.Split(secret, parts, threshold)
	if err != nil {
		return nil, err
	}

	keys := []Part{}
	for k, v := range partsMap {
		keys = append(keys, Part{
			Tag:     k,
			Content: v,
		})
	}

	return keys, nil
}

func Combine(parts []Part) ([]byte, error) {
	combinedMap := map[byte][]byte{}
	for _, p := range parts {
		combinedMap[p.Tag] = p.Content
	}

	res, err := shamir.Combine(combinedMap)
	if err != nil {
		return nil, err
	}

	return res, nil
}
