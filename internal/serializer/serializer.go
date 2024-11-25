package serializer

import (
	"bytes"
	"encoding/gob"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func encode(gl routing.GameLog)  ([]byte, error) {
  var data bytes.Buffer
  enc := gob.NewEncoder(&data)
  err := enc.Encode(gl)
  return data.Bytes(), err
}

func decode(data []byte)  (routing.GameLog, error) {
  dec := gob.NewDecoder(bytes.NewReader(data))
  var log routing.GameLog
  err := dec.Decode(&log)
  return log, err
}

