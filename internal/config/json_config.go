package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	errWrongType = errors.New("wrong type")
	errNotFound  = errors.New("not found")
)

// WrapperOption is used to pass additional options to the wrapper
type WrapperOption func(*Loader)

// New will create a config wrapper based on a JSON file
func newJsonConfig(file string, opts ...WrapperOption) (*Loader, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON config file with err: %s", err)
	}

	blob, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON config file with err: %s", err)
	}

	cfg := &Loader{
		rawConfig: blob,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	cfg, err = newFromBytes(cfg.rawConfig)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// newFromBytes will create a config wrapper based on the supplied JSON	string
func newFromBytes(payload []byte) (*Loader, error) {
	var config map[string]interface{}

	err := json.Unmarshal(payload, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON config with err: %s", err)
	}

	return &Loader{
		jsonConfig: config,
	}, nil
}

// WithSub substitutes a template string such as %%ADDRESS%% with a string
func WithSub(tpl, val string) WrapperOption {
	return func(w *Loader) {
		w.rawConfig = bytes.ReplaceAll(w.rawConfig, []byte(tpl), []byte(val))
	}
}

type Loader struct {
	rawConfig  []byte
	jsonConfig map[string]interface{}
}

func (w *Loader) GetInt(key string) (int, error) {
	val, ok := w.jsonConfig[key]
	if !ok {
		return 0, fmt.Errorf("key: %s - %w", key, errNotFound)
	}

	valAsFloat, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("key: %s, val: %v, type: %T - %w", key, val, val, errWrongType)
	}

	return int(valAsFloat), nil
}

func (w *Loader) GetString(key string) (string, error) {
	val, ok := w.jsonConfig[key]
	if !ok {
		return "", fmt.Errorf("key: %s - %w", key, errNotFound)
	}

	valAsString, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("key: %s, val: %v, type: %T - %w", key, val, val, errWrongType)
	}

	return valAsString, nil
}

func (w *Loader) GetFloat(key string) (float64, error) {
	val, ok := w.jsonConfig[key]
	if !ok {
		return 0, fmt.Errorf("key: %s - %w", key, errNotFound)
	}

	valAsFloat, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("key: %s, val: %v, type: %T - %w", key, val, val, errWrongType)
	}

	return valAsFloat, nil
}

func (w *Loader) GetBool(key string) (bool, error) {
	val, ok := w.jsonConfig[key]
	if !ok {
		return false, fmt.Errorf("key: %s - %w", key, errNotFound)
	}

	valAsBool, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("key: %s, val: %v, type: %T - %w", key, val, val, errWrongType)
	}

	return valAsBool, nil
}

func (w *Loader) GetJSON(key string, destination interface{}) error {
	val, ok := w.jsonConfig[key]
	if !ok {
		return fmt.Errorf("key: %s - %w", key, errNotFound)
	}

	// val is map type
	payload, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("key: %s, val: %v, type: %T - %w - err: %s", key, val, val, errWrongType, err)
	}

	err = json.Unmarshal(payload, destination)
	if err != nil {
		return fmt.Errorf("key: %s, val: %v, type: %T - %w - err: %s", key, val, val, errWrongType, err)
	}

	return nil
}
