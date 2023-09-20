package controller

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/response"
)

func Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "上传失败")
		return
	}
	defer file.Close()

	// 获取上传文件的保存路径，默认是 ./uploads，可以通过path设置上传路径
	uploadPath := c.DefaultQuery("path", "./uploads")

	// 确保保存路径存在，如果不存在则创建
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "上传失败")
		return
	}

	// 指定上传文件的目标路径
	destination := filepath.Join(uploadPath, header.Filename)

	// 创建目标文件
	outFile, err := os.Create(destination)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "上传失败")
		return
	}
	defer outFile.Close()

	// 复制上传的文件到目标文件
	_, err = io.Copy(outFile, file)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "上传失败")
		return
	}

	response.Success(c, nil, "文件上传成功")
}

func Download(c *gin.Context) {
	filename := c.Param("filename")

	// 获取下载文件的路径，默认是 ./uploads,可以通过path设置
	downloadPath := c.DefaultQuery("path", "./uploads")

	// 构建完整的文件路径
	filePath := filepath.Join(downloadPath, filename)

	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		response.Fail(c, gin.H{"error": "文件不存在"}, "文件下载失败")
		return
	}

	// 设置下载文件的响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// 设置文件下载的响应类型
	c.Header("Content-Type", "application/octet-stream")

	// 打开文件并将其内容写入响应体
	file, err := os.Open(filePath)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "文件下载失败")
		return
	}
	defer file.Close()

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		response.Fail(c, gin.H{"error": err.Error()}, "文件下载失败")
		return
	}
}
