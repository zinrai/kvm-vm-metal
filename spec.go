package main

import (
	"fmt"
	"strings"
)

type NICSpec struct {
	Model  string
	Bridge string
}

type NICList []NICSpec

func (n *NICList) String() string {
	return fmt.Sprintf("%v", []NICSpec(*n))
}

func (n *NICList) Set(value string) error {
	kv, err := parseKV(value)
	if err != nil {
		return fmt.Errorf("invalid -nic value %q: %w", value, err)
	}

	spec := NICSpec{Model: "e1000e"}
	for k, v := range kv {
		switch k {
		case "model":
			spec.Model = v
		case "bridge":
			spec.Bridge = v
		default:
			return fmt.Errorf("unknown -nic key: %s", k)
		}
	}

	if spec.Bridge == "" {
		return fmt.Errorf("-nic requires bridge=NAME")
	}

	*n = append(*n, spec)
	return nil
}

type DiskSpec struct {
	Size string
	Bus  string
}

type DiskList []DiskSpec

func (d *DiskList) String() string {
	return fmt.Sprintf("%v", []DiskSpec(*d))
}

func (d *DiskList) Set(value string) error {
	kv, err := parseKV(value)
	if err != nil {
		return fmt.Errorf("invalid -disk value %q: %w", value, err)
	}

	spec := DiskSpec{Bus: "sata"}
	for k, v := range kv {
		switch k {
		case "size":
			spec.Size = v
		case "bus":
			spec.Bus = v
		default:
			return fmt.Errorf("unknown -disk key: %s", k)
		}
	}

	if spec.Size == "" {
		return fmt.Errorf("-disk requires size=SIZE")
	}

	*d = append(*d, spec)
	return nil
}

func parseKV(value string) (map[string]string, error) {
	result := map[string]string{}
	for _, pair := range strings.Split(value, ",") {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("not a key=value pair: %s", pair)
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			return nil, fmt.Errorf("empty key in: %s", pair)
		}
		result[k] = v
	}
	return result, nil
}
