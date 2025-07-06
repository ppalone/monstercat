package monstercat

type IDType string

const (
	CatalogId IDType = "catalogId"
	UUID      IDType = "uuid"
)

type getReleaseOpts struct {
	idType IDType
}

type ReleaseOption func(o *getReleaseOpts)

func newGetReleaseOpts() *getReleaseOpts {
	return &getReleaseOpts{
		idType: CatalogId, // default is catalogId
	}
}

func WithIdType(idType IDType) ReleaseOption {
	return func(o *getReleaseOpts) {
		o.idType = idType
	}
}

func (o *getReleaseOpts) build() map[string]string {
	p := make(map[string]string)
	p["idType"] = string(o.idType)
	return p
}
