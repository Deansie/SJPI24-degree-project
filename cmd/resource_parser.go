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

	spec, ok := doc["spec"].(map[string]interface{})
	if !ok {
		return res
	}

	template, ok := spec["template"].(map[string]interface{})
	if !ok {
		return res
	}
	templateSpec, ok := template["spec"].(map[string]interface{})
	if !ok {
		return res
	}
	containers, ok := templateSpec["containers"].([]interface{})
	if !ok {
		return res
	}

	for _, c := range containers {
		cm, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		container := Container{}

		if name, ok := cm["name"].(string); ok {
			container.Name = name
		}

		if resources, ok := cm["resources"].(map[string]interface{}); ok {
			if req, ok := resources["requests"].(map[string]interface{}); ok {
				container.Requests = map[string]string{}
				for k, v := range req {
					if val, ok := v.(string); ok {
						container.Requests[k] = val
					}
				}
			}

			if lim, ok := resources["limits"].(map[string]interface{}); ok {
				container.Limits = map[string]string{}
				for k, v := range lim {
					if val, ok := v.(string); ok {
						container.Limits[k] = val
					}
				}
			}
		}

		res.Containers = append(res.Containers, container)
	}

	return res
}
