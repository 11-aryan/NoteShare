package notes

func (n Note) GetSearchString() string {
	return n.Title
}

func (n Note) GetID() string {
	return n.Id.Hex()
}
