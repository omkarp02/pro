package types

type ACCESS_TOKEN_PAYLOAD struct {
	ID         string
	ProviderId string
	Role       []string
}

type REFRESH_TOKEN_PAYLOAD struct {
	ID         string
	ProviderId string
	Role       []string
}
