package controllers

import (
	"course/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type PPTFileController struct {
	beego.Controller
}

func (this *PPTFileController) Get() {
	method := this.GetString("method")
	if method == "list" {
		id, err := this.GetInt("id")
		if err != nil {
			id = -1
		}
		name := this.GetString("name")
		course_id, err := this.GetInt("course_id")
		if err != nil {
			course_id = -1
		}
		ppts := models.QueryPPTList(id, course_id, name)
		tmpMapArr := make([]interface{}, len(ppts))
		for i, p := range ppts {
			tmpMap := make(map[string]interface{})
			tmpMap["id"] = p.Id
			tmpMap["name"] = p.Name
			tmpMap["course_id"] = p.CourseId.Id
			tmpMapArr[i] = tmpMap
		}
		bodyJSON := simplejson.New()
		bodyJSON.Set("status", "success")
		bodyJSON.Set("data", tmpMapArr)
		body, _ := bodyJSON.Encode()
		this.Ctx.Output.Body(body)
	} else if method == "getfile" {
		id, err := this.GetInt("id")
		if err != nil {
			this.Abort(models.ErrJson("invalid file id"))
		}
		file, err := models.GetPptFileById(id)
		if err != nil {
			this.Abort(models.ErrJson("invalid file id"))
		}
		this.Ctx.Output.Download(file.FilePath)
	}
}

func (this *PPTFileController) Post() {
	course_id, err := this.GetInt("course_id")
	if err != nil {
		this.Abort(models.ErrJson("must have a course id"))
	}
	name := this.GetString("name")
	var ppt models.PptFile
	course, err := models.GetCourseById(course_id)
	if err != nil {
		this.Abort(models.ErrJson("invalid course id"))
	}
	ppt.CourseId = course
	file, head, err := this.GetFile("file")
	defer file.Close()
	if err != nil {
		this.Abort(models.ErrJson("error when trying to get file"))
	}
	ppt.FilePath = "./upload/" + models.GenerateKey() + "__" + head.Filename
	ppt.Name = name
	id, err := models.AddPptFile(&ppt)
	if err != nil {
		this.Abort(models.ErrJson("add ppt file failed, database error"))
	} else {
		this.SaveToFile("file", ppt.FilePath)
		bodyJSON := simplejson.New()
		bodyJSON.Set("status", "success")
		bodyJSON.Set("id", id)
		body, _ := bodyJSON.Encode()
		this.Ctx.Output.Body(body)
	}
}

func (this *PPTFileController) Delete() {
	id, err := this.GetInt("id")
	if err != nil {
		this.Abort("invalid file id")
	}
	err = models.DeletePptFile(id)
	if err != nil {
		this.Abort(models.ErrJson("invalid file id or database error"))
	}
	this.Ctx.Output.Body(models.SuccessJson())
}