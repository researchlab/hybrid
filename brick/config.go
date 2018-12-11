package brick

import (
	"encoding/json"
	"fmt"

	"github.com/toolkits/file"
)

// Config interface
type Config interface {
	// GetString return string value from config
	GetString(name string, value ...string) string

	// GetBool return bool value from config
	GetBool(name string, value ...bool) bool

	// GetMapString return string value for a map from config
	GetMapString(name string, field string, value ...string) string

	// GetMapBoolean return boolean value for a map from config
	GetMapBool(name string, field string, value ...bool) bool

	// GetMapint return int value for a map from config
	GetMapInt(name string, field string, value ...int) int

	// GetMap return map from config
	GetMap(name string) map[string]interface{}
}

// JSONConfigService load config from a json config file
type JSONConfigService struct {
	path  string
	props map[string]interface{}
}

// NewJSONConfigService create an instance of JSONConfigService
func NewJSONConfigService(path string) *JSONConfigService {
	return &JSONConfigService{path: path}
}

// Init load config file.
func (p *JSONConfigService) Init() error {
	if props, err := p.load(p.path); err != nil {
		return err
	} else {
		p.props = props
		return nil
	}
}

// GetString return string value from config
func (p *JSONConfigService) GetString(name string, value ...string) string {
	if r := p.props[name]; r != nil {
		return r.(string)
	} else if len(value) > 0 {
		return value[0]
	} else {
		return ""
	}
}

// GetBool return bool value from config
func (p *JSONConfigService) GetBool(name string, value ...bool) bool {
	if r := p.props[name]; r != nil {
		return r.(bool)
	} else if len(value) > 0 {
		return value[0]
	} else {
		return false
	}
}

// GetMapString return string value for a map from config
func (p *JSONConfigService) GetMapString(name string, field string, value ...string) string {
	if r := p.props[name]; r != nil {
		if v := r.(map[string]interface{})[field]; v != nil {
			return v.(string)
		}
	}

	if len(value) > 0 {
		return value[0]
	} else {
		return ""
	}
}

func (p *JSONConfigService) GetMapBool(name string, field string, value ...bool) bool {
	if r := p.props[name]; r != nil {
		if v := r.(map[string]interface{})[field]; v != nil {
			return v.(bool)
		}
	}

	if len(value) > 0 {
		return value[0]
	} else {
		return false
	}
}

func (p *JSONConfigService) GetMapInt(name string, field string, value ...int) int {
	if r := p.props[name]; r != nil {
		if o := r.(map[string]interface{})[field]; o != nil {
			switch v := o.(type) {
			case float64:
				return int(v)			
			default:
				return o.(int)
			}
		}
	}

	if len(value) > 0 {
		return value[0]
	} else {
		return 0
	}
}

// GetMap return map from config
func (p *JSONConfigService) GetMap(name string) map[string]interface{} {
	if r := p.props[name]; r != nil {
		return r.(map[string]interface{})
	} else {
		return map[string]interface{}{}
	}
}

func (p *JSONConfigService) load(path string) (map[string]interface{}, error) {
	if !file.IsExist(path) {
		return nil, fmt.Errorf("file: %s isn't exists", path)
	}

	body, err := file.ToTrimString(path)
	if err != nil {
		return nil, err
	}

	props := map[string]interface{}{}
	err = json.Unmarshal([]byte(body), &props)
	return props, err
}
