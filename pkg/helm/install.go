package helm

import (
	"github.com/pkg/errors"
	"github.com/tkeel-io/kit/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func installChart(name, chart, version string, injects ...*chart.Chart) error { // nolint
	installClient := action.NewInstall(defaultCfg)
	valueOpts := &values.Options{}
	installClient.Version = version
	if installClient.Version == "" && installClient.Devel {
		log.Debug("setting version to >0.0.0-0")
		installClient.Version = ">0.0.0-0"
	}
	installClient.ReleaseName = name

	var err error
	cp, err := installClient.ChartPathOptions.LocateChart(chart, env)
	if err != nil {
		err = errors.Wrap(err, "get helm chart path options err")
		return err
	}

	log.Debugf("CHART PATH: %s\n", cp)

	p := getter.All(env)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		err = errors.Wrap(err, "merge some value err")
		return err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		err = errors.Wrap(err, "load chart err")
		return err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return err
	}

	if chartRequested.Metadata.Deprecated {
		log.Warn("This chart is deprecated")
	}

	// Add inject dependencies
	if err := checkInjects(injects, name); err != nil {
		return errors.Wrap(err, "get injects dependency chart err")
	}
	if len(injects) == 0 {
		log.Warn("no component request")
	} else {
		chartRequested.AddDependency(injects...)
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if installClient.DependencyUpdate {
				man := &downloader.Manager{
					Out:              nil,
					ChartPath:        cp,
					Keyring:          installClient.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: env.RepositoryConfig,
					RepositoryCache:  env.RepositoryCache,
					Debug:            env.Debug,
				}
				if err = man.Update(); err != nil {
					return errors.Wrap(err, "helm download manager update err")
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return errors.Wrap(err, "check dependencies err")
			}
		}
	}

	installClient.Namespace = namespace

	if _, err := installClient.Run(chartRequested, vals); err != nil {
		return errors.Wrap(err, "INSTALLATION FAILED")
	}
	return nil
}

func checkInjects(injects []*chart.Chart, pluginName string) error {
	for i := range injects {
		if injects[i] == nil {
			return errors.New("unable dependency chart try to injects")
		}
		injects[i].Values["pluginID"] = pluginName
	}
	return nil
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable.
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
