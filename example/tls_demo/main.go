package main

import (
	"fmt"
	"time"

	tcvectordb "github.com/tencent/vectordatabase-sdk-go/tcvectordb"
)

func main() {
	fmt.Println("=== VectorDB TLS 配置示例（整合版） ===")
	fmt.Println("此文件展示了HTTPS和RPC两种连接方式的TLS配置，使用文件路径方式传参")
	fmt.Println("注意：此文件为演示用途，不包含实际的数据库操作")

	// HTTPS连接示例
	demoHTTPSConnections()

	// RPC连接示例
	demoRPCConnections()

	// 功能总结
	printFeatureSummary()
}

// demoHTTPSConnections 演示HTTPS连接的各种TLS配置
func demoHTTPSConnections() {
	fmt.Println("\n=== HTTPS 连接示例 ===")

	// 示例1: 基本的HTTPS连接，使用根证书文件路径
	fmt.Println("\n1. 基本HTTPS连接（使用CA证书文件路径）:")
	option1 := tcvectordb.ClientOption{
		Timeout: 10 * time.Second,
		CACert:  "path/to/your/ca-cert.pem", // 使用文件路径方式
	}

	client1, err := tcvectordb.NewClient("https://your-vectordb-server.com", "username", "api-key", &option1)
	if err != nil {
		fmt.Printf("   创建HTTPS客户端失败: %v\n", err)
	} else {
		fmt.Println("   ✓ 从文件路径读取CA证书的HTTPS连接创建成功（演示）")
		client1.Close()
	}

	// 示例2: 测试环境跳过证书验证
	fmt.Println("\n2. 测试环境跳过证书验证:")
	option2 := tcvectordb.ClientOption{
		Timeout:            10 * time.Second,
		InsecureSkipVerify: true, // 跳过证书验证（仅限测试环境）
	}

	client2, err := tcvectordb.NewClient("https://your-vectordb-server.com", "username", "api-key", &option2)
	if err != nil {
		fmt.Printf("   创建HTTPS客户端失败: %v\n", err)
	} else {
		fmt.Println("   ✓ 跳过证书验证的HTTPS连接创建成功（演示，仅限测试环境）")
		client2.Close()
	}
}

// demoRPCConnections 演示RPC连接的各种TLS配置
func demoRPCConnections() {
	fmt.Println("\n=== RPC 连接示例 ===")

	// 示例1: 基本的RPC连接，使用HTTP（无TLS）
	fmt.Println("\n1. 基本RPC连接（HTTP）:")
	option1 := tcvectordb.ClientOption{
		Timeout: 10 * time.Second,
	}

	_, err := tcvectordb.NewRpcClientPool("http://your-vectordb-server.com:80", "username", "api-key", &option1)
	if err != nil {
		fmt.Printf("   创建RPC客户端失败: %v\n", err)
	} else {
		fmt.Println("   ✓ 基本RPC连接创建成功（演示）")
	}

	// 示例2: 使用HTTPS和CA证书文件路径的RPC连接
	fmt.Println("\n2. 使用HTTPS和CA证书文件路径的RPC连接:")
	option2 := tcvectordb.ClientOption{
		Timeout: 10 * time.Second,
		CACert:  "path/to/your/ca-cert.pem", // 使用文件路径方式
	}

	_, err = tcvectordb.NewRpcClientPool("https://your-vectordb-server.com:443", "username", "api-key", &option2)
	if err != nil {
		fmt.Printf("   创建RPC客户端失败: %v\n", err)
	} else {
		fmt.Println("   ✓ 从文件路径读取CA证书的RPC连接创建成功（演示）")
	}

	// 示例3: 测试环境跳过证书验证
	fmt.Println("\n3. 测试环境跳过证书验证:")
	option3 := tcvectordb.ClientOption{
		Timeout:            10 * time.Second,
		InsecureSkipVerify: true, // 跳过证书验证（仅限测试环境）
	}

	_, err = tcvectordb.NewRpcClientPool("https://your-vectordb-server.com:443", "username", "api-key", &option3)
	if err != nil {
		fmt.Printf("   创建RPC客户端失败: %v\n", err)
	} else {
		fmt.Println("   ✓ 跳过证书验证的RPC连接创建成功（演示，仅限测试环境）")
	}
}

// printFeatureSummary 打印功能总结
func printFeatureSummary() {
	fmt.Println("\n=== TLS 配置功能总结 ===")
	fmt.Println("\nHTTPS连接功能:")
	fmt.Println("✓ 支持HTTPS连接")
	fmt.Println("✓ 支持CA证书验证（文件路径方式）")
	fmt.Println("✓ 支持SNI设置（默认vdb.tencentcloud.com）")
	fmt.Println("✓ 支持跳过证书验证（测试环境）")
	fmt.Println("✓ 向后兼容HTTP连接")

	fmt.Println("\nRPC连接功能:")
	fmt.Println("✓ 支持HTTP和HTTPS连接")
	fmt.Println("✓ 支持CA证书验证（文件路径方式）")
	fmt.Println("✓ 支持SNI设置（自动处理）")
	fmt.Println("✓ 支持跳过证书验证（测试环境）")
	fmt.Println("✓ 提供完整的数据库操作接口")
	fmt.Println("✓ 高性能的gRPC通信协议")

	fmt.Println("\n=== 使用说明 ===")
	fmt.Println("1. 将 'path/to/your/ca-cert.pem' 替换为实际的CA证书文件路径")
	fmt.Println("2. 将服务器地址、用户名和API密钥替换为实际值")
	fmt.Println("3. 跳过证书验证仅限测试环境使用")
}
