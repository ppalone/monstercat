package monstercat

// Artist.
type Artist struct {
	CatalogRecordID string
	ID              string
	Name            string
	Public          bool
	Role            string
	URI             string
}

// Artist API Response.
type artistAPIResponse struct {
	CatalogRecordID string `json:"CatalogRecordId"`
	ID              string `json:"Id"`
	Name            string `json:"Name"`
	ProfileFileID   string `json:"ProfileFileId"`
	Public          bool   `json:"Public"`
	Role            string `json:"Role"`
	URI             string `json:"URI"`
}

func (r *artistAPIResponse) toArtist() Artist {
	return Artist{
		CatalogRecordID: r.CatalogRecordID,
		ID:              r.ID,
		Name:            r.Name,
		Public:          r.Public,
		Role:            r.Role,
		URI:             r.URI,
	}
}
