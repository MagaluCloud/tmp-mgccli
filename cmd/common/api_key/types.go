package apikey

// ApiKeys representa o retorno da requisição de API Keys
type ApiKeys struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	ApiKey        string  `json:"api_key"`
	Description   string  `json:"description"`
	KeyPairID     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
	Tenant        struct {
		UUID      string `json:"uuid"`
		LegalName string `json:"legal_name"`
	} `json:"tenant"`
	Scopes []struct {
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		ConsentText string `json:"consent_text"`
		Icon        string `json:"icon"`
		APIProduct  struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"api_product"`
	} `json:"scopes"`
}
