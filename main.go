package main

import (
	"fmt"
	"os"

	"aliyun-oss-website-action/config"
	"aliyun-oss-website-action/operation"
	"aliyun-oss-website-action/utils"
)

func main() {
	defer utils.TimeCost()()
	if config.Folder == "/" {
		fmt.Println("You should not upload the root directory, use ./ instead. 通常来说, 你不应该上传根目录, 也许你是要配置 ./")
		os.Exit(1)
	}

	if !config.SkipSetting {
		operation.SetStaticWebsiteConfig()
	} else {
		fmt.Println("skip setting static pages related configuration")
	}

	var incremental *operation.IncrementalConfig
	if config.IsIncremental {
		fmt.Println("---- [incremental] ---->")
		incremental, _ = operation.GetRemoteIncrementalConfig(config.Bucket)
		fmt.Println("<---- [incremental end] ----")
		fmt.Println()
	}
	if !config.IsIncremental || incremental == nil {
		// TODO: delete after upload
		fmt.Println("---- [delete] ---->")
		deleteErrs := operation.DeleteObjects(config.Bucket)
		utils.LogErrors(deleteErrs)
		fmt.Println("<---- [delete end] ----")
		fmt.Println()
	}

	records := utils.WalkDir(config.Folder)

	fmt.Println("---- [upload] ---->")
	uploaded, uploadErrs := operation.UploadObjects(config.Folder, config.Bucket, records, incremental, config.OssDir)
	utils.LogErrors(uploadErrs)
	fmt.Println("<---- [upload end] ----")
	fmt.Println()

	if config.IsIncremental && incremental != nil {
		fmt.Println("---- [delete] ---->")
		deleteErrs := operation.DeleteObjectsIncremental(config.Bucket, incremental)
		utils.LogErrors(deleteErrs)
		fmt.Println("<---- [delete end] ----")
		fmt.Println()
	}

	if config.IsIncremental {
		fmt.Println("---- [incremental] ---->")
		operation.UploadIncrementalConfig(config.Bucket, uploaded)
		fmt.Println("<---- [incremental end] ----")
		fmt.Println()
	}

	if len(uploadErrs) > 0 {
		os.Exit(1)
	}

}
