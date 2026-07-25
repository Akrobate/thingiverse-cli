package thing

type License struct {
	Value       string
	Description string
}

// Licenses: cc, cc-sa, cc-nd, cc-nc-sa, cc-nc-nd, pd0, gpl, lgpl, bsd
var SupportedLicenses = []License{
	{
		Value:       "c",
		Description: "Standard Copyright - All rights reserved",
	},
	{
		Value:       "cc-sa",
		Description: "Creative Commons Attribution-ShareAlike",
	},
	{
		Value:       "cc-nd",
		Description: "Creative Commons Attribution-NoDerivatives",
	},
	{
		Value:       "cc-nc-sa",
		Description: "Creative Commons Attribution-NonCommercial-ShareAlike",
	},
	{
		Value:       "cc-nc-nd",
		Description: "Creative Commons Attribution-NonCommercial-NoDerivatives",
	},
	{
		Value:       "pd0",
		Description: "Public Domain / CC0 Zero Dedication",
	},
	{
		Value:       "gpl",
		Description: "GNU General Public License (Copyleft)",
	},
	{
		Value:       "lgpl",
		Description: "GNU Lesser General Public License",
	},
	{
		Value:       "bsd",
		Description: "BSD Permissive License",
	},
}

func IsValidLicense(value string) bool {
	for _, lic := range SupportedLicenses {
		if lic.Value == value {
			return true
		}
	}
	return false
}
