package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joyent/kosh/v3/conch"
	"github.com/joyent/kosh/v3/logger"
	"github.com/joyent/kosh/v3/tables"
	"github.com/joyent/kosh/v3/template"
)

// Config is the default configuration struct
type Config struct {
	Version string
	GitRev  string

	ConchURL   string
	ConchToken string
	ConchENV   string

	OutputJSON bool

	logger.Logger
}

// NewConfig takes a Version and a GitRev and returns a Config object
func NewConfig(Version, GitRev string) Config {
	return Config{
		Version: Version,
		GitRev:  GitRev,
		Logger:  logger.New(),
	}
}

const configTemplate = `
---
# Config

* Version: {{ .Version }}
* GitRev: {{ .GitRev }}

* ConchENV: {{ .ConchENV }}
* ConchURL: {{ .ConchURL }}
* ConchToken: {{ .ConchToken }}

* OutputJSON: {{ .OutputJSON }}

Logger

* Debug {{ .Logger.LevelDebug  }}
* Info {{ .Logger.LevelInfo  }}
---
`

// String returns a string implementation of the config object
func (c Config) String() string {
	t, err := template.NewTemplate().Parse(configTemplate)
	if err != nil {
		log.Fatal(err)
	}

	buf := &strings.Builder{}
	if e := t.Execute(buf, c); e != nil {
		log.Fatal(e)
	}
	return buf.String()
}

// ConchClient returns a configured client for the Conch API
func (c Config) ConchClient() *conch.Client {
	c.Debug("Creating Conch Client")
	return conch.New(
		conch.API(c.ConchURL),
		conch.AuthToken(c.ConchToken),
		conch.UserAgent(fmt.Sprintf("kosh %s", c.GitRev)),
		conch.Logger(c.Logger),
	)
}

// Renderer returns a function that will render to STDOUT
func (c Config) Renderer() func(interface{}) {
	return c.RenderTo(os.Stdout)
}

func renderJSON(i interface{}) string {
	b, e := json.Marshal(i)
	if e != nil {
		fatal(e)
	}
	return string(b)
}

// RenderTo returns a function tha renders to a given io.Writer based on the
// configuraton and datatype
func (c Config) RenderTo(w io.Writer) func(interface{}) {
	return func(i interface{}) {
		if c.OutputJSON {
			c.Debug("Outputting JSON")
			fmt.Fprintln(w, renderJSON(i))
			return
		}
		switch t := i.(type) {
		case template.Templated:
			s, e := template.Render(t)
			if e != nil {
				fatal(e)
			}
			fmt.Fprintln(w, s)
		case tables.Tabulable:
			fmt.Fprintln(w, tables.Render(t))
		default:
			fmt.Fprintln(w, renderJSON(t))
		}
	}
}
