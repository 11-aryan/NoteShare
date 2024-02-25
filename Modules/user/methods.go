package user

func (u User) GetSearchString() string {
	return u.Name
}

func (u User) GetID() string {
	return u.Id.Hex()
}
