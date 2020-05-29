package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Soar 分析传入的 sql 语句
// 安装 sql
// on windows go to  https://github.com/XiaoMi/soar/releases download soar-windows-amd64
// rename soar-windows-amd64 to soar.exe
// add thr soar.exe working dir into system path
// for example, put soar.exe into d:/soar/soar.exe just add D:/soar into system path
// so that we can use soar in powershell
func Soar() {
	{
		// -query string
		// 待评审的 SQL 或 SQL 文件，如 SQL 中包含特殊字符建议使用文件名。
		// on windows system we can make good use of cmd
		// cmd := exec.Command("cmd", "/C", "soar", "-version")
		cmd := exec.Command("soar", "-query", "D:/soar/car.sql")
		// fmt.Println(cmd.Path, cmd)
		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("in all caps: %q\n", string(out))
	}
	// {
	// 	cmd := exec.Command("soar", "-query", "D:/soar/car.sql")
	// 	// cmd := exec.Command("cmd", "/C", "soar", "-version")
	// 	fmt.Println(cmd.Path, cmd)
	// 	var out bytes.Buffer
	// 	cmd.Stdout = &out
	// 	err := cmd.Run()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("in all caps: %q\n", out.String())
	// }
	// {
	// 	path, err := exec.LookPath("soar")
	// 	if err != nil {
	// 		log.Fatal("installing soar is in your future")
	// 	}
	// 	fmt.Printf("soar is available at %s\n", path)
	// }
}

// LocalSoar use local file to test soar
func LocalSoar() (string, error) {
	// on windows system we can make good use of cmd
	// cmd := exec.Command("cmd", "/C", "soar", "-version") -explain-format json
	cmd := exec.Command("soar", "-query", "D:/soar/car.sql", "-report-type", "json")
	fmt.Println(cmd.Path, cmd)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	fmt.Printf("in all caps: %q\n", string(out))
	return string(out), nil
}

func tempSQLFile(sql string) (string, error) {
	dir := os.TempDir()
	filePath := filepath.Join(dir, "temp.sql")
	outputFile, outputError := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Println("An error occurred with file opening or creation\n", outputError)
		return "", outputError
	}
	defer outputFile.Close()
	outputFile.WriteString(sql)
	return filePath, nil
}

func removeTempSQLFile(filePath string) error {
	return os.Remove(filePath)
}
