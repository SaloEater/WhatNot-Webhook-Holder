package service

import "bytes"

type DigitalOceaner interface {
	SaveLabel(bytes.Buffer, string) (string, error)
}
