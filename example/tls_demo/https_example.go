package main

import (
	"fmt"
	"time"

	tcvectordb "github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func main() {
	fmt.Println("=== VectorDB HTTPS 连接示例 ===")

	// 示例1: 基本的HTTPS连接，使用根证书文件
	fmt.Println("\n1. 基本HTTPS连接:")
	option1 := tcvectordb.ClientOption{
		Timeout: 10 * time.Second,
		CACert:  "path/to/your/ca-cert.pem",
	}

	client1, err := tcvectordb.NewClient("https://your-vectordb-server.com", "username", "api-key", &option1)
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
	} else {
		fmt.Println("✓ 从文件路径读取CA证书的HTTPS连接创建成功")
		client1.Close()
	}

	// 示例2: 使用自定义CA证书
	fmt.Println("\n2. 使用CA证书的HTTPS连接:")
	caCertContent := `-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----`

	option2 := tcvectordb.ClientOption{
		Timeout: 10 * time.Second,
		CACert:  caCertContent, // 直接传入CA证书内容
	}

	client2, err := tcvectordb.NewClient("https://your-vectordb-server.com", "username", "api-key", &option2)
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
	} else {
		fmt.Println("✓ 使用CA证书的HTTPS连接创建成功")
		client2.Close()
	}

	// 示例3: 测试环境跳过证书验证
	fmt.Println("\n4. 测试环境跳过证书验证:")
	option4 := tcvectordb.ClientOption{
		Timeout:            10 * time.Second,
		InsecureSkipVerify: true, // 跳过证书验证（仅限测试环境）
	}

	client4, err := tcvectordb.NewClient("https://your-vectordb-server.com", "username", "api-key", &option4)
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
	} else {
		fmt.Println("✓ 跳过证书验证的HTTPS连接创建成功（仅限测试环境）")
		client4.Close()
	}

	fmt.Println("\n=== 功能总结 ===")
	fmt.Println("✓ 支持HTTPS连接")
	fmt.Println("✓ 支持CA证书验证（直接内容或文件路径）")
	fmt.Println("✓ 支持SNI设置（默认vdb.tencentcloud.com）")
	fmt.Println("✓ 支持跳过证书验证（测试环境）")
	fmt.Println("✓ 向后兼容HTTP连接")
}
