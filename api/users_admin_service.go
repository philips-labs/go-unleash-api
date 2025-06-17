package api

type UserDetails struct {
	Id         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	ImageUrl   string `json:"imageUrl,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	InviteLink string `json:"inviteLink,omitempty"`
	EmailSent  bool   `json:"emailSent,omitempty"`
	RootRole   int    `json:"rootRole,omitempty"`
}

type UsersService struct {
	client *ApiClient
}

func (p *UsersService) SearchUser(query string) (*[]UserDetails, *Response, error) {
	if query == "" {
		return nil, nil, ErrRequiredParam("query")
	}
	req, _ := p.client.newRequest("admin/user-admin/search?q="+query, "GET", nil)

	var users []UserDetails

	resp, err := p.client.do(req, &users)
	if err != nil {
		return nil, resp, err
	}
	return &users, resp, err
}
