package groundwork

import (
	"github.com/profe-ajedrez/sapo/pkg/json/formatter"
)

type Relationable interface {
	Load() error
}

type Relation struct {
	Table    string `default:""`
	Refered  []string
	Refering []string
}

func (bytes *Relation) ToJson() (string, error) {
	return formatter.ToJson(bytes)
}

func (bytes *Relation) ToPrettyJson() (string, error) {
	return formatter.ToPrettyJson(bytes)
}

type BD interface {
	Close() error
	Relations(table string) (Relation, error)
	GetStructure(table string) ([]string, error)
}

type Connection interface {
	Connect(configFilePath string) error
}
