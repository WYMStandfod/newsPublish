package controllers

import (
	"github.com/astaxie/beego"
	"path"
	"time"
	"shanghaiyiqi/models"
	"math"
	"github.com/astaxie/beego/orm"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页
func(this*ArticleController)ShowArticleList(){
	userName:=this.GetSession("userName")
	if userName==nil{
		this.Redirect("/login.html",302)
		return 
	}
	//获取数据
	//高级查询
	//指定表
	o := orm.NewOrm()
	qs := o.QueryTable("Article")//queryseter
	var articles []models.Article
	//_,err := qs.All(&articles)
	//if err != nil{
	//	beego.Info("查询数据错误")
	//}
	typeName:=this.GetString("select")
	//查询总记录数
	var count int64

	//获取总页数
	pageSize := 2

//获取页码
	pageIndex,err:= this.GetInt("pageIndex")
	if err != nil{
		pageIndex = 1
	}

	//获取数据
	//作用就是获取数据库部分数据,第一个参数，获取几条,第二个参数，从那条数据开始获取,返回值还是querySeter
	//起始位置计算
	start := (pageIndex - 1)*pageSize
	if typeName==""{
		count,_ = qs.Count()
	}else{
		count,_ =qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).Count()

	}
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types

    if typeName==""{
		qs.Limit(pageSize,start).All(&articles)

	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__TypeName",typeName).All(&articles)

	}



	//传递数据
	this.Data["userName"]=userName
	this.Data["typeName"]=typeName
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count
	this.Data["articles"] = articles
	this.Layout="layout.html"
	this.TplName = "index.html"
}
//展示添加文章页面
func(this*ArticleController)ShowAddArticle(){
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types

	this.TplName = "add.html"
}
//获取添加文章数据
func(this*ArticleController)HandleAddArticle(){
	//1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	//2校验数据
	if articleName == "" || content == ""{
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}

	//处理文件上传
	file ,head,err:=this.GetFile("uploadname")
	defer file.Close()
	if err != nil{
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}


	//1.文件大小
	if head.Size > 5000000{
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return
	}

	//2.文件格式
	//a.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = "add.html"
		return
	}

	//3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	//存储
	this.SaveToFile("uploadname","./static/img/"+fileName)



	//3.处理数据
	//插入操作
	o := orm.NewOrm()

	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = "/static/img/"+fileName
    typeName:=this.GetString("select")
    var articleType models.ArticleType
    articleType.TypeName=typeName
    //beego.Info(typeName,articleType)
    o.Read(&articleType,"TypeName")
    article.ArticleType=&articleType
	o.Insert(&article)


	//4.返回页面
	this.Redirect("/showArticleList",302)
}
func UploadFile(this*beego.Controller,filePath string)string{
	file ,head,err:=this.GetFile(filePath)
	if head.Filename == ""{
		return "NoImg"
	}
	defer file.Close()
	if err != nil{
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return ""
	}


	//1.文件大小
	if head.Size > 5000000{
		this.Data["errmsg"] = "文件太大，请重新上传"
		this.TplName = "add.html"
		return ""
	}

	//2.文件格式


	//a.jpg
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg"{
		this.Data["errmsg"] = "文件格式错误。请重新上传"
		this.TplName = "add.html"
		return ""
	}

	//3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	//存储
	this.SaveToFile(filePath,"./static/img/"+fileName)
	return "/static/img/"+fileName
}
//展示文章详情页面
func(this*ArticleController)ShowArticleDetail(){
	//获取数据
	id,er:=this.GetInt("articleId")
	//数据校验
	if er != nil{
		beego.Info("传递的链接错误")
	}
	//操作数据
      o:=orm.NewOrm()
      var article models.Article
      article.Id=id
      o.Read(&article)
      o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",id).One(&article)
      article.Acount+=1
      o.Update(&article)

      m2m:=o.QueryM2M(&article,"Users")
      userName:=this.GetSession("userName")
      if userName==nil{
      	this.Redirect("/login",302)
      	return
	  }
	  var user models.User
	  user.Name=userName.(string)
	  o.Read(&user,"Name")
      m2m.Add(user)
	//返回视图页面
     var users []models.User
     o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	userLayout:=this.GetSession("userName")
     this.Data["userName"]=userLayout.(string)
     this.Data["pageType"]="文章详情"
    this.Data["users"]=users
	this.Data["article"] = article
	this.Layout="layout.html"
	this.TplName = "content.html"
}
func(this*ArticleController)ShowUpdateArticle(){
	id,err:=this.GetInt("articleId")
	if err!=nil{
		beego.Info("请求错误")
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	o.Read(&article)
	this.Data["pageType"]="更新文章"
	this.Layout= "layout.html"
	this.Data["article"]=article
	this.TplName="update.html"
}
func(this*ArticleController)HandleUpdateArticle(){
	id,err:=this.GetInt("articleId")
	articleName:=this.GetString("articleName")
	content:=this.GetString("content")
	filePath:=UploadFile(&this.Controller,"uploadname")
	if err!=nil||articleName==""||content==""||filePath==""{
		beego.Info("请求路径错误")
		return
	}
	o:=orm.NewOrm()
	var article models.Article
	article.Id=id
	err=o.Read(&article)
	if err!=nil{
		beego.Info("数据不存在")
		return
	}
	article.ArtiName=articleName
	article.Acontent=content
	if filePath!="NoImg"{
		article.Aimg=filePath
	}

	o.Update(&article)
	this.Redirect("/article/showArticleList",302)
}
func (this*ArticleController)DeleteArticle()  {
	 id,err:=this.GetInt("articleId")
	 if err!=nil{
	 	beego.Info("删除文章请求路径错误")
	 }
	 o:=orm.NewOrm()
	 var article models.Article
	 article.Id=id
	 o.Delete(&article)
	 this.Redirect("/showArticleList",302)
}
func (this*ArticleController)ShowAddType(){
	o:=orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types
	this.TplName="addType.html"
}
func(this*ArticleController)HandleAddType(){
	typeName:=this.GetString("typeName")
	if typeName==""{
		beego.Info("信息不完整,请重新输入")
	}
	o:=orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName=typeName
	o.Insert(&articleType)
	this.Redirect("/addType",302)
}
func(this*ArticleController)DeleteType(){
   id,err:=this.GetInt("id")
   if err!=nil{
   	beego.Error("删除类型错误",err)
   	return
   }
   o:=orm.NewOrm()
   var articleType models.ArticleType
   articleType.Id=id
   o.Delete(&articleType)
   this.Redirect("/article/addType",302)
}