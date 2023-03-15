package main

import (
	// 导入Chainmaker2.0的SDK包
	sdk "chainmaker.org/chainmaker/sdk-go/v2"

	// 导入其他需要的包
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

const sdk_config = "./testdata/sdk_config.yml"

// 定义一个结构体来存储JSON数据
type Data struct {
	ContractAddress string            `json:"contractName"` // 合约地址
	FunctionName    string            `json:"contractFunc"` // 合约函数名
	Parameters      map[string]string `json:"params"`       // 合约函数参数
}

// 定义一个处理HTTP请求的函数
func handler(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("读取请求体出错:", err)
		return
	}
	defer r.Body.Close()

	// 将JSON数据解析为Data结构体
	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("解析JSON数据出错:", err)
		return
	}

	// fmt.Println("收到的JSON数据:", data)

	// 使用默认配置创建一个客户端实例
	c, _ := sdk.NewChainClient(sdk.WithConfPath(sdk_config))

	var c_params []*common.KeyValuePair
	for k, v := range data.Parameters {
		pair := &common.KeyValuePair{Key: k, Value: []byte(v)}
		c_params = append(c_params, pair)
	}

	// 调用合约函数并获取结果为字节切片
	result, err := c.InvokeContract(data.ContractAddress, data.FunctionName, "", c_params, -1, true)
	if err != nil {
		fmt.Println("调用合约函数出错:", err)
		return
	}

	// fmt.Println("合约函数结果:", result)

	// 将结果写回响应写入器为十六进制字符串
	w.Write(result.ContractResult.Result)
}

func main() {
	// 在8080端口监听和服务
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
