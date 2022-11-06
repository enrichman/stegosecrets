package stego

import (
	"encoding/base64"

	shamir "github.com/corvus-ch/shamir"
	"github.com/pkg/errors"
)

type Part struct {
	Version   byte
	Parts     byte
	Threshold byte
	Tag       byte
	Content   []byte
}

func NewPartFromContent(content []byte) (Part, error) {
	if len(content) < 5 {
		return Part{}, errors.New("invalid part: ot enough content bytes")
	}

	return NewPart(
		content[0],
		content[1],
		content[2],
		content[3],
		content[4:],
	), nil
}

func NewPart(version, parts, threshold, tag byte, content []byte) Part {
	return Part{
		Version:   version,
		Parts:     parts,
		Threshold: threshold,
		Tag:       tag,
		Content:   content,
	}
}

func (p Part) Bytes() []byte {
	return append([]byte{
		p.Version,
		p.Parts,
		p.Threshold,
		p.Tag,
	}, p.Content...)
}

func (p Part) Base64() string {
	return base64.StdEncoding.EncodeToString(p.Bytes())
}

func Split(secret []byte, parts, threshold uint8) ([]Part, error) {
	partsMap, err := shamir.Split(secret, int(parts), int(threshold))
	if err != nil {
		return nil, errors.Wrap(err, "failed splitting secret")
	}

	keys := []Part{}
	for k, v := range partsMap {
		keys = append(keys, NewPart(
			'1',
			parts,
			threshold,
			k,
			v,
		))
	}

	return keys, nil
}

func Combine(parts []Part) ([]byte, error) {
	combinedMap := map[byte][]byte{}
	for _, p := range parts {
		combinedMap[p.Tag] = p.Content
	}

	if len(combinedMap) < int(parts[0].Threshold) {
		return nil, errors.Errorf(
			"not enough parts provided: parts %d, threshold %d",
			len(combinedMap), parts[0].Threshold,
		)
	}

	res, err := shamir.Combine(combinedMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed combining secret")
	}

	return res, nil
}
