package service

import (
	"bufio"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/http"
	"score-system/dao"
	"score-system/model"
	"strconv"
	"strings"
)

type Service struct {
	dao *dao.Dao
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		dao: dao.NewDao(db),
	}
}

func (s *Service) UpdateFile(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.String(http.StatusBadRequest, "get from file err: %+v", err)
		return
	}
	files := form.File["files"]
	log.Printf("上传文件个数：%d \n", len(files))

	achievements := make([]*model.Achievement, 0)
	errChan := make(chan error)
	dchan := make(chan *model.Achievement, len(files))
	quitChan := make(chan bool, len(files))
	doNum := 0

	for _, f := range files {
		go readFile(f, errChan, dchan, quitChan)
	}

loop:
	for {
		select {
		case e := <-errChan:
			log.Println(e)
			ctx.String(http.StatusInternalServerError, "red data error: %+v", e)
			return
		case d, ok := <-dchan:
			if !ok {
				break loop
			}
			achievements = append(achievements, d)
		case <-quitChan:
			doNum++
			if doNum == len(files) {
				close(dchan)
			}
		}
	}

	log.Printf("成绩数组：len=%d, data=%+v \n", len(achievements), achievements)
	// 将成绩写入数据库
	err = s.dao.BatchInsertAchievement(achievements)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, "成绩上传成功")
}

func readFile(fileHeader *multipart.FileHeader, errChan chan error, dChan chan *model.Achievement, quitChan chan bool) {
	file, err := fileHeader.Open()
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		str := scan.Text()
		log.Println(str)
		data := strings.Split(str, "、")
		if len(data) != 4 {
			errChan <- errors.New("成绩格式不正确")
			return
		}

		class, err := strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			errChan <- err
			return
		}
		name := data[1]
		score, err := strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			errChan <- err
			return
		}
		category := data[3]

		achievement := &model.Achievement{
			Class:   int32(class),
			Name:    name,
			Score:   int32(score),
			Subject: category,
		}

		dChan <- achievement
	}
	quitChan <- true
}

func (s *Service) GetAchievement(ctx *gin.Context) {
	var params QueryCondition
	if err := ctx.ShouldBind(&params); err != nil {
		log.Printf("[GetAchievement] 请求参数解析失败：%+v", err)
		ctx.String(http.StatusBadRequest, "请求参数有误")
	}

	condition := &dao.Condition{
		Start:   params.Start,
		Limit:   params.Limit,
		Class:   params.Class,
		Subject: params.Subject,
		Name:    params.Name,
	}
	if params.OrderType != 0 {
		if params.OrderType == 1 {
			condition.OrderType = dao.ByASC
		} else {
			condition.OrderType = dao.ByDESC
		}
	}
	if params.OrderFiled != 0 {
		switch params.OrderFiled {
		case Subject:
			condition.OrderFiled = dao.BySubject
		case Class:
			condition.OrderFiled = dao.ByClass
		case Score:
			condition.OrderFiled = dao.ByScore
		}
	}

	count, data, err := s.dao.GetAchievementOrderByFiled(condition)
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusInternalServerError, "系统错误：%+v", err)
		return
	}
	res := QueryRes{
		Count:        count,
		Achievements: data,
	}
	log.Println(data)
	ctx.JSON(http.StatusOK, res)
}
