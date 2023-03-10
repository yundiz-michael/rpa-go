package rpc

import (
	"bytes"
	"image"
	"image/jpeg"
	"merkaba/common"
)

/**
	See视觉解析，使用方法：
    如果有图标，可以再server.yaml文件中定义，如下，就定义京麦眼睛的图标，并给他分配了一个文本名称
	icons:
   		-jm_eye: "icons/jm_eye.jpg"
	params的参数描述
    1) 如果包含icons=["jm_eye"],标识plato会返回图标jm_eye的坐标
    2) 如果detectText=True|False,标识是否需要识别文本坐标
    基于此，在go的客户端就能实现基于视觉的自动化操作
	select("jm_eye").click();
    select("用户名称",left|right|top|bottom ,offsetX)
*/

func SeeImage(ip string, params map[string]any, jpegImg image.Image) (map[string]any, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, jpegImg, nil)
	if err != nil {
		return nil, err
	}
	return See(ip, params, buf.Bytes())
}

func See(ip string, params map[string]any, bytes []byte) (map[string]any, error) {
	client := UsePaasClient(ip, "IMAGE", common.LocalIP)
	result, err := client.Invoke("See", params, bytes)
	client.Dispose()
	if err != nil {
		common.LoggerStd.Error(err.Error())
		return nil, err
	}
	return result.(map[string]any), nil
}
