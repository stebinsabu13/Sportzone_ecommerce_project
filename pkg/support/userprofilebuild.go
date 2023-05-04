package support

import "github.com/stebinsabu13/ecommerce-api/pkg/utils"

func BuildProfile(details utils.ResponseUsers, address []utils.Address) utils.ResponseUserDetails {
	return utils.ResponseUserDetails{
		FirstName: details.FirstName,
		LastName:  details.LastName,
		Email:     details.Email,
		MobileNum: details.MobileNum,
		Address:   address,
	}
}
