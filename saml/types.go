package saml

type SAMLInfo struct {
	SessionName     string
	SessionDuration int
	Url             string
	Accounts        []*Account
	RawSAML         string
}

type Role struct {
	Name      string
	Arn       string
	AccountID string
	Url       string
}

type Account struct {
	Alias string
	ID    string
	Url   string
	Roles []*Role
}

type AWSSAMLService interface {
	SAMLRequestURL() (string, error)
	ParseSAMLResponse(samlResponse string) (*SAMLInfo, error)
}
