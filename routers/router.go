// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"bethapi/controllers"

	"github.com/astaxie/beego"
)

func Router() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/current",&controllers.CurrentBlockController{}),
		beego.NSRouter("/getblockbynum/:num",&controllers.BlockController{},"get:GetBlock"),
		beego.NSRouter("/getlatestblocks",&controllers.BlockController{},"get:GetLatestBlocks"),
		beego.NSRouter("/getblocksbypage/:page",&controllers.BlockController{},"get:GetBlockByPage"),
	)
	beego.AddNamespace(ns)
}
