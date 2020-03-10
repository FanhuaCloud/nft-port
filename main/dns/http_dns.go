package dns

import (
	"github.com/asmcos/requests"
)

func Resolve(domain string) (string, error) {
	req := requests.Requests()
	resp, err := req.Get("http://119.29.29.29/d?dn=" + domain)
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}
