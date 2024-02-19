package blogService

import (
	"learn_redis/backend/common"
	"learn_redis/backend/config"
	"learn_redis/backend/model"
	"learn_redis/backend/service/userService"
)

func GetHotBlogs(page_no int) common.Result {
	var blog []model.Blog
	config.MySQL.Order("liked DESC").Limit(common.PAGE_SIZE).Offset((page_no - 1) * common.PAGE_SIZE).Find(&blog)
	for i := 0; i < len(blog); i++ {
		var user model.User
		config.MySQL.Select("nick_name", "icon").Where("id = ?", blog[i].UserID).Take(&user)
		blog[i].Icon = user.Icon
		blog[i].Name = user.NickName
	}
	return common.Result{Success: true, Data: blog}
}

func GetUserBlogs(user userService.UserDTO, page_no int) common.Result {
	var blog []model.Blog
	config.MySQL.Where("user_id = ?", user.ID).Order("create_time DESC").Limit(common.PAGE_SIZE).Offset((page_no - 1) * common.PAGE_SIZE).Find(&blog)
	return common.Result{Success: true, Data: blog}
}
