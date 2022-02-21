package mcheads

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

const McHeadsAPIendpoint = "https://mc-heads.net"

func GetAvaterURI(identifier string) string {
	return McHeadsAPIendpoint + "/avatar/" + identifier
}

func GetAvater(identifier string) (io.ReadCloser, error) {
	resp, err := http.Get(GetAvaterURI(identifier))
	if err != nil {
		return nil, errors.Wrap(err, "Get")
	}

	return resp.Body, nil
}
