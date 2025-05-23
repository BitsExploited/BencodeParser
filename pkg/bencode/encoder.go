package bencode

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

// BencodeEncoder interface for custom types that want to define their bencode representation.
type Bencoder interface {
	Bencode() (string, error)
}

func Encode(data interface{}) (string, error) {
	var buf bytes.Buffer
	err := encodeToBuffer(data, &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func encodeToBuffer(data interface{}, writer io.Writer) error {
	switch v := data.(type) {
	case Bencoder: // Custom bencoding
		s, err := v.Bencode()
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(s))
		return err
	case string:
		_, err := fmt.Fprintf(writer, "%d:%s", len(v), v)
		return err
	case int:
		_, err := fmt.Fprintf(writer, "i%de", v)
		return err
	case int64:
		_, err := fmt.Fprintf(writer, "i%de", v)
		return err
	case []interface{}: // List
		if _, err := writer.Write([]byte{'l'}); err != nil {
			return err
		}
		for _, item := range v {
			if err := encodeToBuffer(item, writer); err != nil {
				return err
			}
		}
		if _, err := writer.Write([]byte{'e'}); err != nil {
			return err
		}
		return nil
	case map[string]interface{}: // Dictionary
		if _, err := writer.Write([]byte{'d'}); err != nil {
			return err
		}
		// Keys must be sorted for canonical bencoding
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			// Encode key (must be string)
			if err := encodeToBuffer(k, writer); err != nil {
				return fmt.Errorf("failed to encode dictionary key '%s': %w", k, err)
			}
			// Encode value
			if err := encodeToBuffer(v[k], writer); err != nil {
				return fmt.Errorf("failed to encode dictionary value for key '%s': %w", k, err)
			}
		}
		if _, err := writer.Write([]byte{'e'}); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported type for bencoding: %T", v)
	}
}
