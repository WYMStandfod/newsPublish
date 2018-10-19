package routers

import (
	"shanghaiyiqi/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,Filter)
    beego.Router("/", &controllers.MainController{},"get:ShowGet;post:Post")

    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandlePost")

    beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
   //文章列表页访问
    beego.Router("/article/showArticleList",&controllers.ArticleController{},"get:ShowArticleList")
    //添加文章
    beego.Router("/article/addArticle",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
    //显示文章详情

    beego.Router("/article/showArticleDetail",&controllers.ArticleController{},"get:ShowArticleDetail")
    beego.Router("/article/updateArticle",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
	beego.Router("/article/deleteArticle",&controllers.ArticleController{},"get:DeleteArticle")
    beego.Router("/article/addType",&controllers.ArticleController{},"get:ShowAddType;post:HandleAddType")
    beego.Router("/article/logout",&controllers.UserController{},"get:Logout")
    beego.Router("/article/deleteType",&controllers.ArticleController{},"get:DeleteType")
    }
var Filter = func(ctx *context.Context){
	userName:=ctx.Input.Session("userName")
	if userName==nil{
		ctx.Redirect(302,"/login")
		return
	}
}