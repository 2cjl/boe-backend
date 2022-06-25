package service

import (
	jwtx "boe-backend/internal/util/jwt"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	var registerForm jwtx.RegisterForm
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request parameter",
		})
		return
	}

	//isExist, err := db.IsExistPhone(registerForm.Phone)
	//if err != nil {
	//	log.Println(err)
	//	c.JSON(500, gin.H{
	//		"error": "Server internal error",
	//	})
	//	return
	//}
	//if isExist {
	//	c.JSON(400, gin.H{
	//		"error": "phone already exists",
	//	})
	//	return
	//}
	//
	//id, err := db.Register(registerForm.Username, registerForm.Phone, registerForm.Passwd, registerForm.IdNumber, registerForm.WorkStatus, registerForm.Age)
	//if err != nil {
	//	log.Println(err)
	//	c.JSON(500, gin.H{
	//		"error": "Server internal error",
	//	})
	//	return
	//}

	//fmt.Printf("id (%v, %T)\n", id, id)

	c.JSON(200, gin.H{
		//"token": jwtx.GenerateToken(id),
	})
}
