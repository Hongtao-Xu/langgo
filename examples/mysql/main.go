package main

import (
	"langgo"
	"langgo/components/mysql"
	"langgo/core"
	"langgo/core/log"
)

func main() {
	langgo.Run(&mysql.Instance{})
	//开发环境下，自动建表
	if core.EnvName == core.Development {
		mysql.Main().AutoMigrate(&Account{})
	}
	acc := Account{Name: "langgo"}
	mysql.Main().Create(&acc)
	acc.Name = "famingjia"
	mysql.Main().Save(&acc)
	newAcc := Account{}
	mysql.Main().First(&newAcc, "id=?", acc.ID)
	log.Logger("app", "main").Info().Interface("newAcc", newAcc).Send()
	mysql.Main().Unscoped().Delete(&Account{}, newAcc.ID)
}
