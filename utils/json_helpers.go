package utils

import (
    "bytes"
    "encoding/json"
    "io"
)

// DeserializerJSON decodes a JSON body into the given type.
func DeserializerJSON[T any](body io.Reader) (T, error) {
    var result T
    err := json.NewDecoder(body).Decode(&result)
    return result, err
}

// SerializerJSON encodes the given data into a JSON byte buffer.
func SerializerJSON[T any](data T) (*bytes.Buffer, error) {
    var buffer bytes.Buffer
    err := json.NewEncoder(&buffer).Encode(data)
    return &buffer, err
}