package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type PptFile struct {
	Id       int     `orm:"column(id);auto"`
	CourseId *Course `orm:"column(course_id);rel(fk)"`
	Name     string  `orm:"column(name);size(40);null"`
	FilePath string  `orm:"column(file_path);size(200);null"`
}

func (t *PptFile) TableName() string {
	return "ppt_file"
}

func init() {
	orm.RegisterModel(new(PptFile))
}

// AddPptFile insert a new PptFile into database and returns
// last inserted Id on success.
func AddPptFile(m *PptFile) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPptFileById retrieves PptFile by Id. Returns error if
// Id doesn't exist
func GetPptFileById(id int) (v *PptFile, err error) {
	o := orm.NewOrm()
	v = &PptFile{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllPptFile retrieves all PptFile matches certain condition. Returns empty list if
// no records exist
func GetAllPptFile(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(PptFile))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PptFile
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdatePptFile updates PptFile by Id and returns error if
// the record to be updated doesn't exist
func UpdatePptFileById(m *PptFile) (err error) {
	o := orm.NewOrm()
	v := PptFile{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePptFile deletes PptFile by Id and returns error if
// the record to be deleted doesn't exist
func DeletePptFile(id int) (err error) {
	o := orm.NewOrm()
	v := PptFile{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&PptFile{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func QueryPPTList(id int, course_id int, name string) []PptFile {
	var ppts []PptFile
	o := orm.NewOrm()
	qs := o.QueryTable(new(PptFile))
	cond := orm.NewCondition()
	if id != -1 {
		cond = cond.And("Id", id)
	}
	if course_id != -1 {
		course, err := GetCourseById(course_id)
		if err != nil {
			return ppts
		}
		cond = cond.And("CourseId", course)
	}
	qs.SetCond(cond).Filter("Name__contains", name).All(&ppts)
	return ppts
}
