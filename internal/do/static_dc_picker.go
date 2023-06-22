package do

type StaticDigitalOceanDataCenterPicker struct {
	data_center string
}

func NewStaticDigitalOceanDataCenterPicker(data_center string) *StaticDigitalOceanDataCenterPicker {
	return &StaticDigitalOceanDataCenterPicker{
		data_center: data_center,
	}
}

func (g *StaticDigitalOceanDataCenterPicker) Pick() string {
	return g.data_center
}
