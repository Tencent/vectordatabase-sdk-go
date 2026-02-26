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
MIIDADCCAeigAwIBAgIUPRUK8GRdhEDgNPwv/75FWvaELTkwDQYJKoZIhvcNAQEL
BQAwGDEWMBQGA1UEAxMNa2hhb3Mtcm9vdC1jYTAeFw0yNTEyMDIxMTU0MDVaFw00
NTExMjcxMTU0MDVaMBgxFjAUBgNVBAMTDWtoYW9zLXJvb3QtY2EwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDFBzP49atCkDLFTyg6wKdFNXLF4ijCIH2w
NoyNugqx0ZeZ2GsIkWl73o0vTgrjGkfi/Ze4A2TIFUs5iHTbPRoUnQaLDj+/oNy0
pORuENWRQKyHxukvizma9APUSi/yZcX2lAyHD2IMog5hxHGC4jJKqGirC0b9bQCv
XsOwxpOY0/u6K4KH7msGUlSGE5dFz45fZkVqxPSKkflfKBCV/EonW6EqKn4wi/wk
SIEXzkuG4FpTggoK9eBMqabNa2ZC+YwxGS0T86MYNypY9uFzDWHl4Qd44iBNifH/
B6EoshwmjZ8bBU+gHeaqIdh3/NnVsE7nHjdv2xgLPCm0ZWV03ParAgMBAAGjQjBA
MA4GA1UdDwEB/wQEAwICBDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTucNAD
2Uf0lRId/6cDHgkvyxHz0jANBgkqhkiG9w0BAQsFAAOCAQEAMvUNIsnwQPVg5fSQ
RDlb0PcRnereHLuRJhKPx5ufhqIse8oO/yT1ELPqiLmfzRGVRyL9jjkip64xzgu+
lJBTswk5Vksf/NZkmXF9A7fMrvkZELFlpxG5BE8vq+HXz/OqNPqXVKbETI3J5Vsv
umqxCIGtOINzmj+ccmy+azmKWayvOeWZx5837Y4KwJYHH41B/+sqtCayEvu+ZlN5
qVfhKTRDSjJyE96f1bjEJKdZ8xlvRiG2N5+9rogvUN0gS1JnwBqxqRei0Vu6tFOr
ezV6WGgWYZ5iuYCuX+NWsdU8HLk+lUYPPxkIONEH+TOqBl4zlTPAQ6voSYGqFnMT
CLmkVA==
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
