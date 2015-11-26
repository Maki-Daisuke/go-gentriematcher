package triematcher

import (
	"io"
	"text/template"
)

var templateFile = template.Must(template.New("file").Parse(`
// DO NOT EDIT!
// Code generated by go-triematcher <https://github.com/Maki-Daisuke/go-triematcher>
// DO NOT EDIT!

package {{ .PackageName }}

func Match{{ .TagName }}String(str string) bool {
  return Match{{ .TagName}}(([]byte)(str))
}

func Match{{ .TagName }}(bytes []byte) bool {
  defer func(){
    recover() // Must be "index out of range" error for string.
              // Ignore and return false.
  }()

  i := 0

{{ $start := .Start }}
  STATE_{{ $start.Id }}:
{{ if .IsGoal }}
    return true
{{ else }}
    switch bytes[i] {
  {{ range $key, $next := $start.Nexts }}
    case {{ printf "%q" $key }}:
      i++
      goto STATE_{{ $next.Id }}
  {{ end }}
    default:
      i++
      goto STATE_{{ $start.Id }}
    }
{{ end }}

{{ range .States }}
  STATE_{{ .Id }}:
  {{ if .IsGoal }}
      return true
  {{ else }}
    switch bytes[i] {
    {{ range $key, $next := .Nexts }}
    case {{ printf "%q" $key }}:
      i++
      goto STATE_{{ $next.Id }}
    {{ end }}
    default:
      goto STATE_{{ $start.Id }}
    }
  {{ end }}
{{ end }}
}
`))

func generate(w io.Writer, pkg_name, tag_name string, st *state) error {
	states := listStates(st)
	err := templateFile.Execute(w, map[string]interface{}{
		"PackageName": pkg_name,
		"TagName":     tag_name,
		"Start":       states[0],
		"States":      states[1:],
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Unlike allStates, this does not traverse goal state.
func listStates(start *state) []*state {
	marked := map[int]bool{}
	states := []*state{}

	var traverse func(*state)
	traverse = func(s *state) {
		if marked[s.Id] {
			return
		}
		states = append(states, s)
		marked[s.Id] = true
		if s.IsGoal {
			return
		}
		for _, next := range s.Nexts {
			traverse(next)
		}
	}
	traverse(start)

	return states
}
