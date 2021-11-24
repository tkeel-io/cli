package helm

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/pkg/errors"
	"github.com/tkeel-io/cli/pkg/output"
	"helm.sh/helm/v3/cmd/helm/search"
	"io"
)

func searchAll() (output.Writer, error) {
	index, err := buildIndex()
	if err != nil {
		return nil, errors.Wrap(err, "build index failed")
	}

	var res []*search.Result
	res = index.All()
	search.SortScore(res)
	data, err := applyConstraint(defaultSelectVersion, res)
	if err != nil {
		return nil, errors.Wrap(err, "apply constraint failed")
	}

	return &repoSearchWriter{data, 50}, err
}

type repoChartElement struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

type repoSearchWriter struct {
	results     []*search.Result
	columnWidth uint
}

func (r *repoSearchWriter) WriteTable(out io.Writer) error {
	if len(r.results) == 0 {
		_, err := out.Write([]byte("No results found\n"))
		if err != nil {
			return fmt.Errorf("unable to write results: %s", err)
		}
		return nil
	}
	table := uitable.New()
	table.MaxColWidth = r.columnWidth
	table.AddRow("NAME", "CHART VERSION", "APP VERSION", "DESCRIPTION")
	for _, r := range r.results {
		table.AddRow(r.Name, r.Chart.Version, r.Chart.AppVersion, r.Chart.Description)
	}
	return output.EncodeTable(out, table)
}

func (r *repoSearchWriter) WriteJSON(out io.Writer) error {
	return r.encodeByFormat(out, output.JSON)
}

func (r *repoSearchWriter) WriteYAML(out io.Writer) error {
	return r.encodeByFormat(out, output.YAML)
}

func (r *repoSearchWriter) encodeByFormat(out io.Writer, format output.Format) error {
	// Initialize the array so no results returns an empty array instead of null
	chartList := make([]repoChartElement, 0, len(r.results))

	for _, r := range r.results {
		chartList = append(chartList, repoChartElement{r.Name, r.Chart.Version, r.Chart.AppVersion, r.Chart.Description})
	}

	switch format {
	case output.JSON:
		return output.EncodeJSON(out, chartList)
	case output.YAML:
		return output.EncodeYAML(out, chartList)
	}

	// Because this is a non-exported function and only called internally by
	// WriteJSON and WriteYAML, we shouldn't get invalid types
	return nil
}
