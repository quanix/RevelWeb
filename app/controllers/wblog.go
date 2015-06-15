package controllers
import (
	"RevelWeb/app/models"
	"github.com/revel/revel"
	"strings"
	"fmt"
)

type WBlog struct {
	App
}


//提交博客
func (c WBlog) Putup(blog *models.Blog) revel.Result {
	blog.Title = strings.TrimSpace(blog.Title)
	blog.Email = strings.TrimSpace(blog.Email)
	blog.Subject = strings.TrimSpace(blog.Subject)
	blog.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		fmt.Println(c.Validation)
		return c.Redirect(App.WBlog)
	}
	dao, err := models.NewDao()
	if err != nil {
		c.Response.Status = 500
		return c.RenderError(err)
	}

	defer dao.Close()
	err = dao.CreateBlog(blog)
	if err != nil {
		c.Response.Status = 500
		return c.RenderError(err)
	}
	return c.Redirect(App.Index)
}