package controllers
import (
	"RevelWeb/app/models"
	"github.com/revel/revel"
	"strings"
)

type WComment struct {
	App
}

func (c WComment) Docomment(id string, rcnt int, comment *models.Comment) revel.Result {
	if len(id) == 0 {
		return c.Redirect(App.Index)
	}

	dao, err := models.NewDao()
	if err != nil {
		c.Response.Status = 500
		return c.Redirect(App.Index)
	}

	defer dao.Close()
	blog := dao.FindBlogById(id)
	if blog == nil {
		return c.Redirect(App.Index)
	}

	comment.BlogId = blog.Id
	comment.Content = strings.TrimSpace(comment.Content)
	comment.Email = strings.TrimSpace(comment.Email)
	comment.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Errs:The email and the content should not be null,or the maxsize of email is 50.")
		return c.Redirect("/bloginfor/%s/%d", id, rcnt)
	}

	err = dao.InsertComment(comment)
	if err != nil {
		c.Response.Status = 500
		return c.RenderError(err)
	}

	blog.CommentCnt ++
	dao.UpdateBlogById(id, blog)
	return c.Redirect("/bloginfor/%s/%d", id, rcnt)

}

