package models
import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/revel/revel"
)

type Blog struct {
	Id           bson.ObjectId
	Email        string
	CDate        time.Time
	Title        string
	ShortTitle   string
	Subject      string
	ShortSubject string
	CommentCnt   int
	ReadCnt      int
	Year         int
}

func (blog *Blog) Validate(v *revel.Validation) {
	v.Check(blog.Title,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{200})

	v.Check(blog.Email,
		revel.Required{},
		revel.MaxSize{50})

	v.Email(blog.Email)

	v.Check(blog.Subject,
		revel.Required{},
		revel.MinSize{1})
}

func (dao *Dao) CreateBlog(blog *Blog) error {
	blogCollection := dao.session.DB(DbName).C(BlogCollection)
	blog.Id = bson.NewObjectId()
	blog.CDate = time.Now()
	blog.Year = blog.CDate.Year()

	_, err := blogCollection.Upsert(bson.M{"_id":blog.Id}, blog)
	if err != nil {
		revel.WARN.Printf("Unable to save blog : %v error %v", blog, err)
	}
	return err
}


//获取前50条博客
func (dao *Dao) FindBlog() []Blog {
	blogCollection := dao.session.DB(DbName).C(BlogCollection)
	blogs := []Blog{}
	query := blogCollection.Find(bson.M{}).Sort("-cdate").Limit(50)
	query.All(&blogs)
	return blogs
}


func (blog *Blog) GetShortTitle() string {
	if len(blog.Title)>35 {
		return blog.Title[:35]
	}
	return blog.Title
}
func (blog *Blog) GetShortContent() string {
	if len(blog.Subject)>200 {
		return blog.Subject[:200]
	}
	return blog.Subject
}


func (dao *Dao) FindBlogById(id string) *Blog {
	blogCollection := dao.session.DB(DbName).C(BlogCollection)
	blog := new(Blog)
	query := blogCollection.Find(bson.M{"id":bson.ObjectIdHex(id)})
	query.One(blog)
	return blog
}

func (dao *Dao) UpdateBlogById(id string, blog *Blog) {
	blogCollection := dao.session.DB(DbName).C(BlogCollection)
	err := blogCollection.Update(bson.M{"id": bson.ObjectIdHex(id)}, blog)
	if err != nil {
		revel.WARN.Printf("Unable to update blog: %v error %v", blog, err)
	}
}