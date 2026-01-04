package cmd

import (
	"io"

	"gopkg.in/yaml.v3"
)

func ParseResources(r io.Reader) ([]Resource, error) {
	var resources []Resource

	decoder := yaml.NewDecoder(r)

	for {
		var doc map[string]interface{}
		if err := decoder.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		res := extractResource(doc)
		if res.Kind != "" {
			resources = append(resources, res)
		}
	}

	return resources, nil
}

func extractResource(doc map[string]interface{}) Resource {
	res := Resource{}

	if kind, ok := doc["kind"].(string); ok {
		res.Kind = kind
	}

	if meta, ok := doc["metadata"].(map[string]interface{}); ok {
		if name, ok := meta["name"].(string); ok {
			res.Name = name
		}
		if ns, ok := meta["namespace"].(string); ok {
			res.Namespace = ns
		}
		if labels, ok := meta["labels"].(map[string]interface{}); ok {
			res.Labels = make(map[string]string)
			for k, v := range labels {
				if val, ok := v.(string); ok {
					res.Labels[k] = val
				}
			}
		}
	}

	return res
}
