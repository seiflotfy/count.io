package sketches

import (
	"github.com/retailnext/hllpp"

	"datamodel"
	pb "datamodel/protobuf"
	"utils"
)

// HLLPPSketch is the toplevel sketch to control the HLL implementation
type HLLPPSketch struct {
	*datamodel.Info
	impl *hllpp.HLLPP
}

// NewHLLPPSketch ...
func NewHLLPPSketch(info *datamodel.Info) (*HLLPPSketch, error) {
	d := HLLPPSketch{info, hllpp.New()}
	return &d, nil
}

// Add ...
func (d *HLLPPSketch) Add(values [][]byte) (bool, error) {
	dict := make(map[string]uint)
	for _, v := range values {
		dict[string(v)]++
	}
	for v := range dict {
		d.impl.Add([]byte(v))
	}
	return true, nil
}

// Get ...
func (d *HLLPPSketch) Get(interface{}) (interface{}, error) {
	return &pb.CardinalityResult{
		Cardinality: utils.Int64p(int64(d.impl.Count())),
	}, nil
}
