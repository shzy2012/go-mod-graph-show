package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/* example usage:
// 1.生成依赖关系
go mod graph > deps.txt

// 2.生成 graph.dot
echo "digraph G {" > graph.dot
cat deps.txt | awk '{print "  \"" $1 "\" -> \"" $2 "\";"}' >> graph.dot
echo "}" >> graph.dot

// 3.生成graph图片
dot -Tpng graph.dot -o graph.png
dot -Tpng graph.dot -o graph.svg
*/

const format = `digraph G {
%s
	}`

func main() {
	lines := strings.Builder{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			fmt.Printf("%s -> %s\n", parts[0], parts[1])
			lines.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", parts[0], parts[1]))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
	}

	os.WriteFile("graph.dot", fmt.Appendf(nil, format, lines.String()), 0644)

	// dot -Tpng graph.dot -o graph.png
	// dot -Tsvg graph.dot -o graph.svg

	outputFormat := "png"
	// 检查输出格式是否有效
	if len(os.Args) > 1 {
		outputFormat = os.Args[1]
	}

	// 检查 dot 命令是否存在
	_, err := exec.LookPath("dot")
	if err != nil {
		fmt.Printf("未找到 'dot' 命令，请确保 Graphviz 已安装并在系统 PATH 中")
		return
	}

	outputFile := fmt.Sprintf("graph.%s", outputFormat)
	// 定义命令和参数
	cmd := exec.Command("dot", fmt.Sprintf("-T%s", outputFormat), "graph.dot", "-o", outputFile)

	// 运行命令
	err = cmd.Run()
	if err != nil {
		fmt.Printf("命令执行失败: %s\n", err.Error())
		return
	}

	fmt.Printf("graph文件生成成功: %s\n", outputFile)
}
