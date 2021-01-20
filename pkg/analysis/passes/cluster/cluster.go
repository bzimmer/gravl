package cluster

import (
	"flag"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"github.com/rs/zerolog/log"

	"github.com/bzimmer/gravl/pkg/analysis"
)

const doc = `clusters returns the activities clustered by (distance, elevation) dimensions`

type kMeans struct {
	Clusters  int
	Threshold float64
}

type Cluster struct {
	Center     []float64            `json:"center"`
	Activities []*analysis.Activity `json:"activities"`
}

// results converts the internal cluster struct to one suitable for use externally
func results(c clusters.Clusters) []*Cluster {
	var res []*Cluster
	if c == nil {
		return res
	}
	for i := 0; i < len(c); i++ {
		x := &Cluster{Center: c[i].Center, Activities: make([]*analysis.Activity, 0)}
		for j := 0; j < len(c[i].Observations); j++ {
			obs := c[i].Observations[j].(*observation)
			x.Activities = append(x.Activities, obs.Activity)
		}
		res = append(res, x)
	}
	return res
}

type observation struct {
	Activity *analysis.Activity   `json:"activity"`
	Coords   clusters.Coordinates `json:"coordinates"`
}

func (obs *observation) Coordinates() clusters.Coordinates {
	return obs.Coords
}

func (obs *observation) Distance(point clusters.Coordinates) float64 {
	return obs.Coords.Distance(point)
}

func (k *kMeans) run(ctx *analysis.Context, pass *analysis.Pass) (interface{}, error) {
	if len(pass.Activities) < k.Clusters {
		log.Warn().Int("n", len(pass.Activities)).Int("clusters", k.Clusters).Msg("too few activities")
		return results(nil), nil
	}
	// For each activity, create a synthetic coordinate from the distance and elevation
	//  scaled between 0.0 and 1.0.

	// 1. Create two slices, one for distance and one for elevation
	// 2. Find the max of each slice
	var dmax, emax float64
	var dsts, elvs []float64
	for i := 0; i < len(pass.Activities); i++ {
		act := pass.Activities[i]
		dst := act.Distance.Meters()
		dsts = append(dsts, dst)
		if dst > dmax {
			dmax = dst
		}
		elv := act.ElevationGain.Meters()
		elvs = append(elvs, elv)
		if elv > emax {
			emax = elv
		}
	}

	// 3. Divide each element by the max value
	var d clusters.Observations
	for i := 0; i < len(pass.Activities); i++ {
		d = append(d, &observation{
			Activity: analysis.ToActivity(pass.Activities[i], ctx.Units),
			Coords:   clusters.Coordinates{dsts[i] / dmax, elvs[i] / emax},
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
	return results(clusters), nil
}

func New() *analysis.Analyzer {
	k := &kMeans{Threshold: 0.01, Clusters: 4}
	fs := flag.NewFlagSet("cluster", flag.ExitOnError)
	fs.IntVar(&k.Clusters, "clusters", k.Clusters, "number of clusters")
	fs.Float64Var(&k.Threshold, "threshold", k.Threshold, `threshold (in percent between 0.0 and 0.1) aborts processing
if less than n% of data points shifted clusters in the last iteration`)
	return &analysis.Analyzer{
		Name:  fs.Name(),
		Doc:   doc,
		Flags: fs,
		Run:   k.run,
	}
}
