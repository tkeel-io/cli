package helm

import (
	"io"

	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/tkeel-io/cli/pkg/output"
	"helm.sh/helm/v3/pkg/repo"
)

func listRepo() (*repoListWriter, error) {
	f, err := repo.LoadFile(env.RepositoryConfig)
	if isNotExist(err) || len(f.Repositories) == 0 {
		return nil, errors.New("no repositories to show")
	}

	return &repoListWriter{f.Repositories}, nil
}

type repoListWriter struct {
	repos []*repo.Entry
}

func (r *repoListWriter) WriteJSON(out io.Writer) error {
	return r.encodeByFormat(out, output.JSON)
}

func (r *repoListWriter) WriteYAML(out io.Writer) error {
	return r.encodeByFormat(out, output.YAML)
}

func (r repoListWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	table.AddRow("NAME", "URL")
	for _, re := range r.repos {
		table.AddRow(re.Name, re.URL)
	}
	return output.EncodeTable(out, table)
}

func (r *repoListWriter) encodeByFormat(out io.Writer, format output.Format) error {
	// Initialize the array so no results returns an empty array instead of null
	repolist := make([]repositoryElement, 0, len(r.repos))

	for _, re := range r.repos {
		repolist = append(repolist, repositoryElement{Name: re.Name, URL: re.URL})
	}

	switch format {
	case output.JSON:
		return output.EncodeJSON(out, repolist)
	case output.YAML:
		return output.EncodeYAML(out, repolist)
	case output.TABLE:
		return r.WriteTable(out)
	}

	// Because this is a non-exported function and only called internally by
	// WriteJSON and WriteYAML, we shouldn't get invalid types
	return nil
}

type repositoryElement struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
