package controller

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"gitee.com/openeuler/PilotGo/cmd/server/app/config"
	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gitee.com/openeuler/PilotGo/sdk/response"
	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")

	if contentType == "multipart/form-data" {
		// 直接读取request body内容
		bodyBuf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("没获取到request body: %s", err.Error())
			response.Fail(c, gin.H{"error": err.Error()}, "获取文件request body失败")
			return
		}

		parsedURL, err := url.Parse(c.Request.RequestURI)
		if err != nil {
			logger.Error("解析 URL 错误:%v", err.Error())
			response.Fail(c, gin.H{"error": err.Error()}, "解析URI失败")
			return
		}
		filename := parsedURL.Query().Get("filename")

		uploadPath := c.DefaultQuery("path", config.OptionsConfig.Storage.Path) // 获取上传文件的保存路径，可以通过path设置上传路径
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {            // 确保保存路径存在，如果不存在则创建
			response.Fail(c, gin.H{"error": err.Error()}, "创建保存路径失败")
			return
		}
		destination := filepath.Join(uploadPath, filename) // 上传文件的目标路径

		outFile, err := os.Create(destination) // 创建并打开目标文件
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "创建目标文件失败")
			return
		}
		defer outFile.Close()

		_, err = outFile.Write(bodyBuf) // 将请求体中的二进制数据写入目标文件
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "文件保存失败")
			return
		}

	} else {
		// api测试工具及网页前端上传文件
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "未获取到文件")
			return
		}
		defer file.Close()

		uploadPath := c.DefaultQuery("path", config.OptionsConfig.Storage.Path) // 获取上传文件的保存路径,可以通过path设置上传路径

		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil { // 确保保存路径存在，如果不存在则创建
			response.Fail(c, gin.H{"error": err.Error()}, "保存路径创建失败")
			return
		}

		destination := filepath.Join(uploadPath, header.Filename) // 上传文件的目标路径

		outFile, err := os.Create(destination) // 创建目标文件
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "目标文件创建失败")
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file) // 复制上传的文件到目标文件
		if err != nil {
			response.Fail(c, gin.H{"error": err.Error()}, "文件保存失败")
			return
		}
	}
	response.Success(c, nil, "文件上传成功")
}

func Download(c *gin.Context) {
	filename := c.Param("filename")

	// 获取下载文件的路径，可以通过path设置
	downloadPath := c.DefaultQuery("path", config.OptionsConfig.Storage.Path)

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
