package data

type Brand struct {
	Name    string            `bson:"_id"`
	Logo    map[string]string `bson:"logo"`
	Website string            `bson:"website"`
}

func NewBrand(Name string, idTrendee int64, logo map[string]string, website string) *Brand {
	return &Brand{
		Name:    Name,
		Logo:    logo,
		Website: website,
	}
}

func (b *Brand) Public() PublicBrand {
	return PublicBrand{
		Name:    b.Name,
		Logo:    b.Logo,
		Website: b.Website,
	}
}

type PublicBrand struct {
	Name    string            `json:"name"`
	Logo    map[string]string `json:"image"`
	Website string            `json:"website"`
}
