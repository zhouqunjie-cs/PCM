package tenanter

type AccessKeyTenant struct {
	name   string
	id     string
	secret string
	url    string
	token  string
}

func NewTenantWithAccessKey(name, accessKeyId, accessKeySecret, url, token string) Tenanter {
	return &AccessKeyTenant{
		name:   name,
		id:     accessKeyId,
		secret: accessKeySecret,
		url:    url,
		token:  token,
	}
}

func (tenant *AccessKeyTenant) AccountName() string {
	return tenant.name
}

func (tenant *AccessKeyTenant) Clone() Tenanter {
	return &AccessKeyTenant{
		id:     tenant.id,
		secret: tenant.secret,
		name:   tenant.name,
		url:    tenant.url,
		token:  tenant.token,
	}
}

func (tenant *AccessKeyTenant) GetId() string {
	return tenant.id
}

func (tenant *AccessKeyTenant) GetSecret() string {
	return tenant.secret
}

func (tenant *AccessKeyTenant) GetUrl() string {
	return tenant.url
}

func (tenant *AccessKeyTenant) GetToken() string {
	return tenant.token
}
