package handler

import (
	"io/ioutil"
	"mime/multipart"

	"ginfra/errcode"
	"ginfra/log"
	"ginfra/protocol"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//UploadRequest
type UploadRequest struct {
	ActivityID     string                `form:"ActivityID", binding:"required,min=1"`
	InvitationFile *multipart.FileHeader `form:"InvitationFile"`
}

//UploadResponse
type UploadResponse struct {
}

//Upload
func Upload(c *gin.Context) {
	var req UploadRequest
	err := c.ShouldBind(&req)
	if err != nil {
		log.WithGinContext(c).Error(err.Error(), zap.String("error", errcode.ErrInvalidParam))
		protocol.SetErrResponse(c, protocol.ErrCodeInvalidParameter)
		return
	}

	if req.InvitationFile != nil {
		content, err := ReadUploadedFile(req.InvitationFile)
		if err != nil {
			log.WithGinContext(c).Error(err.Error(), zap.String("error", errcode.ErrCodeInternalError))
			protocol.SetErrResponse(c, errcode.NewCustomError(errcode.ErrCodeInternalError, "读取文件失败"))
			return
		}
		log.WithGinContext(c).Debug(string(content))
	}

	//b, _ := json.Marshal(req.InvitationList)
	//log.WithGinContext(c).Debug(string(b))
	data := &UploadResponse{}
	protocol.SetResponse(c, data)
}

// ReadUploadedFile get content form file.
func ReadUploadedFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	body, err := ioutil.ReadAll(src)
	return body, err
}
