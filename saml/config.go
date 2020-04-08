package saml

import "sort"

type Config struct {
	SpUrl                 string            `mapstructure:"sp_url"`
	AccountAliases        map[string]string `mapstructure:"account_aliases"`
	DefaultUrl            string            `mapstructure:"default_url"`
	AccountUrls           map[string]string `mapstructure:"account_urls"`
	DevAccountUrls        map[string]string `mapstructure:"dev_account_urls"`
	IdpMetadata           string            `mapstructure:"idp_metadata"`
	HideUnaliasedAccounts bool              `mapstructure:"hide_unaliased_accounts"`
}

func (c *Config) UpdateMetaData(samlInfo *SAMLInfo) {
	// set aliases and urls on data
	for i, account := range samlInfo.Accounts {
		if alias, ok := c.AccountAliases[account.ID]; ok {
			samlInfo.Accounts[i].Alias = alias
		}

		var accountUrl string
		if url, ok := c.AccountUrls[account.ID]; ok {
			accountUrl = url
		} else {
			accountUrl = c.DefaultUrl
		}

		samlInfo.Accounts[i].Url = accountUrl
		for _, role := range samlInfo.Accounts[i].Roles {
			role.Url = accountUrl
		}
	}

	// filter accounts only by relay state / c
	var filteredAccounts []*Account
	for _, account := range samlInfo.Accounts {
		// if set - hide all unaliased accounts
		if c.HideUnaliasedAccounts && account.Alias == "" {
			continue
		}
		filteredAccounts = append(filteredAccounts, account)
	}

	// sort by alias
	sort.Slice(filteredAccounts, func(i, j int) bool {
		return filteredAccounts[i].Alias > filteredAccounts[j].Alias
	})
	samlInfo.Accounts = filteredAccounts
}
