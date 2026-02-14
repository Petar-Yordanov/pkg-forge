package manifest

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func (sl *StepLike) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("step must be a mapping")
	}

	for i := 0; i+1 < len(value.Content); i += 2 {
		k := value.Content[i]
		v := value.Content[i+1]
		if k.Kind == yaml.ScalarNode && k.Value == "select" {
			var sn SelectNode
			if err := v.Decode(&sn); err != nil {
				return err
			}
			sl.Select = &sn
			sl.Step = nil
			return nil
		}
	}

	var s Step
	if err := value.Decode(&s); err != nil {
		return err
	}
	sl.Step = &s
	sl.Select = nil
	return nil
}
