package handlers

import (
	"fmt"

	"github.com/Micah-Shallom/stage-two/utils"
	"github.com/gin-gonic/gin"
)

func (h *Handlers) GetUserHandler(c *gin.Context) {

	userID := c.Param("id")
	authUserID, exists := c.Get("UserID")

	if !exists {
		utils.BadRequestResponse(c, "Client error", 400, nil)
		return
	}

	//retrieve user from the database
	user, err := h.App.Models.Users.GetByID(userID)
	fmt.Println("user",user)
	if err != nil {
		utils.BadRequestResponse(c, "Client error", 400, err)
		return
	}

	//check if the authenticated user is requesting their own data and if he is not check the userID he is trying to query belongs to an organization he belongs to
	authenticated := utils.IsAuthenticated(user.UserID, authUserID)
	if !authenticated {

		//if the userID is not mine then i check if the user is a member of any of the organizations i belong to or created

		//Fetch organizations the I (the authenticted user) belongs to or has created
		orgs, err := h.App.Models.Organisations.GetByUserID(authUserID.(string))
		if err != nil {
			utils.BadRequestResponse(c, "Client error", 400, err)
			return
		}

		//check if the requested user is a member of any of the organizations
		//if not, return a 403 forbidden error
		userInOrg := false
		for _, org := range orgs {
			for _, u := range org.Users {
				if u.UserID == user.UserID {
					userInOrg = true
					break
				}
			}
			if userInOrg {
				//break out of the loop if user is found in any of the organizations
				break
			}
		}

		if !userInOrg {
			utils.BadRequestResponse(c, "Client error", 400, err)
			return
		}

		//if the user is a member of any of the organizations i belong to or created then i can proceed to get the user data
		utils.SendUserResponse(c, user)
	}
	//if the userID is mine then i can proceed to get the user data
	utils.SendUserResponse(c, user)
}
