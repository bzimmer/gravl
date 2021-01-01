package cluster

import (
	"context"
	"flag"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const Doc = `clusters returns the activities clustered by (distance, elevation) dimensions`

type KMeans struct {
	Clusters  int
	Threshold float64
}

type Observation struct {
	Activity *analysis.Activity   `json:"activity"`
	Coords   clusters.Coordinates `json:"coordinates"`
}

func (obs *Observation) Coordinates() clusters.Coordinates {
	return obs.Coords
}

func (obs *Observation) Distance(point clusters.Coordinates) float64 {
	return obs.Coords.Distance(point)
}

func (k *KMeans) Run(ctx context.Context, pass *analysis.Pass) (interface{}, error) {
	// For each activity, create a synthetic coordinate from the distance and elevation
	//  scaled between 0.0 and 1.0.

	// 1. Create two slices, one for distance and one for elevation
	// 2. Find the max of each slice
	var dm, em float64
	var dsts, elvs []float64
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]

		dst := act.Distance.Meters()
		dsts = append(dsts, dst)
		if dst > dm {
			dm = dst
		}

		elv := act.ElevationGain.Meters()
		elvs = append(elvs, elv)
		if elv > em {
			em = elv
		}
	}

	// 3. Divide each element by the max slice
	var d clusters.Observations
	for i := 0; i < len(pass.Activities); i++ {
		d = append(d, &Observation{
			Activity: analysis.ToActivity(pass.Activities[i], pass.Units),
			Coords:   clusters.Coordinates{dsts[i] / dm, elvs[i] / em},
		})
	}

	// 4. Partition
	km, err := kmeans.NewWithOptions(k.Threshold, nil)
	if err != nil {
		return nil, err
	}
	clusters, err := km.Partition(d, k.Clusters)
	if err != nil {
		return nil, err
	}
	return clusters, nil
}

func New() *analysis.Analyzer {
	k := &KMeans{
		Threshold: 0.01,
		Clusters:  4,
	}
	fs := flag.NewFlagSet("kmeans", flag.ExitOnError)
	fs.IntVar(&k.Clusters, "clusters", k.Clusters, "number of clusters")
	fs.Float64Var(&k.Threshold, "threshold", k.Threshold, `threshold (in percent between 0.0 and 0.1) aborts processing
if less than n% of data points shifted clusters in the last iteration`)
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   Doc,
		Flags: fs,
		Run:   k.Run,
	}
}
