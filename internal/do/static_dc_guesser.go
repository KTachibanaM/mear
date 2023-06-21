package do

type StaticDigitalOceanDataCenterGuesser struct {
	data_center string
}

func NewStaticDigitalOceanDataCenterGuesser(data_center string) *StaticDigitalOceanDataCenterGuesser {
	return &StaticDigitalOceanDataCenterGuesser{
		data_center: data_center,
	}
}

func (g *StaticDigitalOceanDataCenterGuesser) Guess() string {
	return g.data_center
}
