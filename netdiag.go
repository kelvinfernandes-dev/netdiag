package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Configurações básicas
const VERSION = "2.0.0"

// Cores ANSI
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Flags
var (
	flagPing        = flag.String("ping", "", "Testar ping/latência em um host (ex: google.com)")
	flagSpeed       = flag.Bool("speed", false, "Testar velocidade da internet (download e upload)")
	flagInterfaces  = flag.Bool("interfaces", false, "Listar interfaces de rede")
	flagPort        = flag.String("port", "", "Testar uma porta específica (ex: 80 ou localhost:8080)")
	flagDNS         = flag.String("dns", "", "Testar resolução DNS (ex: google.com)")
	flagIP          = flag.Bool("ip", false, "Mostrar IP público e privado")
	flagTrace       = flag.String("trace", "", "Executar traceroute para um host")
	flagAll         = flag.Bool("all", false, "Executar todos os testes de diagnóstico")
	flagJSON        = flag.Bool("json", false, "Saída em formato JSON")
	flagHelp        = flag.Bool("help", false, "Mostrar ajuda")
	flagVersion     = flag.Bool("version", false, "Mostrar versão")
	flagCount       = flag.Int("count", 4, "Número de pings a executar (padrão: 4)")
	flagInteractive = flag.Bool("i", false, "Modo interativo")
)

type DiagResult struct {
	Test      string      `json:"test"`
	Status    string      `json:"status"`
	Details   interface{} `json:"details,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func main() {
	flag.Parse()

	// Se não tiver flags, entrar no modo interativo
	if len(os.Args) == 1 || *flagInteractive {
		runInteractiveMode()
		return
	}

	if *flagVersion {
		printVersion()
		return
	}

	if *flagHelp {
		printHelp()
		waitForExit()
		return
	}

	var results []DiagResult

	// Executar todos os testes se -all for especificado
	if *flagAll {
		printHeader("EXECUTANDO DIAGNÓSTICO COMPLETO DE REDE")
		results = append(results, runAllTests()...)
	} else {
		// Executar ações baseadas nas flags
		if *flagInterfaces {
			results = append(results, testInterfaces())
		}

		if *flagIP {
			results = append(results, testPublicIP())
		}

		if *flagPing != "" {
			results = append(results, testPing(*flagPing, *flagCount))
		}

		if *flagSpeed {
			results = append(results, testDownloadSpeed(), testUploadSpeed())
		}

		if *flagPort != "" {
			results = append(results, testPortConnection(*flagPort))
		}

		if *flagDNS != "" {
			results = append(results, testDNSResolution(*flagDNS))
		}

		if *flagTrace != "" {
			results = append(results, testTraceroute(*flagTrace))
		}
	}

	// Exibir resultados
	if *flagJSON {
		outputJSON(results)
	} else {
		if len(results) == 0 {
			printWarning("Nenhum teste executado. Use -help para ver opções.")
		} else {
			printResults(results)
		}
		waitForExit()
	}
}

func runInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	// Limpar tela
	clearScreen()

	for {
		printInteractiveMenu()

		fmt.Print("\nEscolha uma opção: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		fmt.Println()

		switch input {
		case "1":
			// Interfaces
			result := testInterfaces()
			printResults([]DiagResult{result})

		case "2":
			// IP Público
			result := testPublicIP()
			printResults([]DiagResult{result})

		case "3":
			// Ping
			fmt.Print("Digite o host (ex: google.com): ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			if host != "" {
				result := testPing(host, 4)
				printResults([]DiagResult{result})
			} else {
				printError("Host não pode estar vazio")
			}

		case "4":
			// DNS
			fmt.Print("Digite o host (ex: google.com): ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			if host != "" {
				result := testDNSResolution(host)
				printResults([]DiagResult{result})
			} else {
				printError("Host não pode estar vazio")
			}

		case "5":
			// Velocidade
			printInfo("Testando velocidade de download...")
			downloadResult := testDownloadSpeed()
			printResults([]DiagResult{downloadResult})

			printInfo("Testando velocidade de upload...")
			uploadResult := testUploadSpeed()
			printResults([]DiagResult{uploadResult})

		case "6":
			// Testar Porta
			fmt.Print("Digite a porta ou host:porta (ex: 80 ou localhost:8080): ")
			port, _ := reader.ReadString('\n')
			port = strings.TrimSpace(port)
			if port != "" {
				result := testPortConnection(port)
				printResults([]DiagResult{result})
			} else {
				printError("Porta não pode estar vazia")
			}

		case "7":
			// Traceroute
			fmt.Print("Digite o host (ex: google.com): ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			if host != "" {
				printInfo("Executando traceroute (isso pode demorar)...")
				result := testTraceroute(host)
				printResults([]DiagResult{result})
			} else {
				printError("Host não pode estar vazio")
			}

		case "8":
			// Diagnóstico Completo
			printHeader("EXECUTANDO DIAGNÓSTICO COMPLETO DE REDE")
			results := runAllTests()
			printResults(results)

		case "9":
			// Ajuda
			printHelp()

		case "0", "sair", "exit", "q":
			printSuccess("Saindo... Até logo!")
			return

		default:
			printWarning("Opção inválida! Digite um número de 0 a 9.")
		}

		fmt.Print("\nPressione ENTER para continuar...")
		reader.ReadString('\n')
		clearScreen()
	}
}

func printInteractiveMenu() {
	fmt.Printf(`%s╔════════════════════════════════════════════════════╗
║      NetDiag v%s - Menu Interativo         ║
╚════════════════════════════════════════════════════╝%s

%s[1]%s Listar Interfaces de Rede
%s[2]%s Mostrar IP Público
%s[3]%s Testar Ping/Latência
%s[4]%s Testar Resolução DNS
%s[5]%s Testar Velocidade (Download + Upload)
%s[6]%s Testar Porta Específica
%s[7]%s Executar Traceroute
%s[8]%s Diagnóstico Completo (todos os testes)
%s[9]%s Mostrar Ajuda
%s[0]%s Sair
`, ColorBold+ColorCyan, VERSION, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorGreen, ColorReset,
		ColorYellow, ColorReset,
		ColorRed, ColorReset)
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printHeader(text string) {
	if !*flagJSON {
		fmt.Printf("\n%s%s%s\n", ColorBold+ColorCyan, text, ColorReset)
		fmt.Println(strings.Repeat("=", len(text)))
	}
}

func printSuccess(text string) {
	if !*flagJSON {
		fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, text)
	}
}

func printError(text string) {
	if !*flagJSON {
		fmt.Printf("%s✗%s %s\n", ColorRed, ColorReset, text)
	}
}

func printWarning(text string) {
	if !*flagJSON {
		fmt.Printf("%s⚠%s %s\n", ColorYellow, ColorReset, text)
	}
}

func printInfo(text string) {
	if !*flagJSON {
		fmt.Printf("%sℹ%s %s\n", ColorBlue, ColorReset, text)
	}
}

func printVersion() {
	fmt.Printf("%sNetDiag v%s%s\n", ColorBold+ColorCyan, VERSION, ColorReset)
	fmt.Printf("Sistema: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Go version: %s\n", runtime.Version())
}

func printHelp() {
	fmt.Printf(`%s%sNetDiag - Ferramenta de Diagnóstico de Rede%s
Versão: %s

%sUSO:%s
  netdiag              (Modo interativo - menu)
  netdiag [OPÇÕES]     (Modo linha de comando)

%sOPÇÕES:%s
  -i                 Forçar modo interativo
  -interfaces        Listar interfaces de rede ativas
  -ip                Mostrar IP público e IPs locais
  -ping <host>       Testar ping/latência (ICMP ou TCP)
  -count <n>         Número de pings (padrão: 4)
  -speed             Testar velocidade de download e upload
  -port <porta>      Testar se uma porta está aberta (ex: 80 ou localhost:8080)
  -dns <host>        Testar resolução DNS
  -trace <host>      Executar traceroute para um host
  -all               Executar todos os testes de diagnóstico
  -json              Saída em formato JSON
  -version           Mostrar versão do programa
  -help              Mostrar esta ajuda

%sEXEMPLOS:%s
  netdiag                          (menu interativo)
  netdiag -interfaces
  netdiag -ping google.com
  netdiag -ping google.com -count 10
  netdiag -speed
  netdiag -port 80
  netdiag -port localhost:8080
  netdiag -dns google.com
  netdiag -trace google.com
  netdiag -all
  netdiag -ip -ping google.com -json

%sNOTAS:%s
  • Dê dois cliques no executável para abrir o menu interativo
  • O teste de ping tenta usar ICMP (requer privilégios) e fallback para TCP
  • O teste de velocidade baixa ~10MB para medição mais precisa
  • Traceroute pode demorar alguns segundos
  • Use -json para integração com outras ferramentas

`, ColorBold+ColorCyan, ColorReset, ColorReset, VERSION,
		ColorBold+ColorYellow, ColorReset,
		ColorBold+ColorYellow, ColorReset,
		ColorBold+ColorYellow, ColorReset,
		ColorBold+ColorYellow, ColorReset)
}

func waitForExit() {
	if runtime.GOOS == "windows" && !*flagJSON {
		fmt.Printf("\n%sPressione ENTER para sair...%s", ColorYellow, ColorReset)
		fmt.Scanln()
	}
}

// FUNÇÕES DE TESTE
func runAllTests() []DiagResult {
	var results []DiagResult

	printInfo("1/6 Testando interfaces de rede...")
	results = append(results, testInterfaces())

	printInfo("2/6 Obtendo IP público...")
	results = append(results, testPublicIP())

	printInfo("3/6 Testando conectividade (google.com)...")
	results = append(results, testPing("google.com", 3))

	printInfo("4/6 Testando resolução DNS (google.com)...")
	results = append(results, testDNSResolution("google.com"))

	printInfo("5/6 Testando velocidade de download...")
	results = append(results, testDownloadSpeed())

	printInfo("6/6 Testando velocidade de upload...")
	results = append(results, testUploadSpeed())

	return results
}

func testInterfaces() DiagResult {
	result := DiagResult{
		Test:      "interfaces",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		result.Status = "error"
		result.Error = err.Error()
		return result
	}

	var activeInterfaces []map[string]string

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					activeInterfaces = append(activeInterfaces, map[string]string{
						"name": iface.Name,
						"ip":   ipnet.IP.String(),
						"mask": ipnet.Mask.String(),
					})
				}
			}
		}
	}

	if len(activeInterfaces) == 0 {
		result.Status = "warning"
		result.Details = "Nenhuma interface ativa encontrada"
	} else {
		result.Status = "success"
		result.Details = activeInterfaces
	}

	return result
}

func testPublicIP() DiagResult {
	result := DiagResult{
		Test:      "public_ip",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha ao obter IP público: %v", err)
		return result
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha ao ler resposta: %v", err)
		return result
	}

	result.Status = "success"
	result.Details = map[string]string{
		"ip": strings.TrimSpace(string(ip)),
	}

	return result
}

func testPing(host string, count int) DiagResult {
	result := DiagResult{
		Test:      "ping",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Tentar ping ICMP primeiro (requer privilégios)
	if runtime.GOOS != "windows" {
		if icmpResult, success := tryICMPPing(host, count); success {
			result.Status = "success"
			result.Details = icmpResult
			return result
		}
	}

	// Fallback para TCP ping
	target := host
	if !strings.Contains(host, ":") {
		target = host + ":80"
	}

	var latencies []float64
	var successful int

	for i := 0; i < count; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, 3*time.Second)
		latency := time.Since(start).Milliseconds()

		if err == nil {
			conn.Close()
			latencies = append(latencies, float64(latency))
			successful++
		}

		if i < count-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	if successful == 0 {
		result.Status = "error"
		result.Error = fmt.Sprintf("Host %s não alcançável", host)
		return result
	}

	// Calcular estatísticas
	var sum, min, max float64
	min = latencies[0]
	max = latencies[0]

	for _, lat := range latencies {
		sum += lat
		if lat < min {
			min = lat
		}
		if lat > max {
			max = lat
		}
	}

	avg := sum / float64(len(latencies))
	loss := float64(count-successful) / float64(count) * 100

	result.Status = "success"
	result.Details = map[string]interface{}{
		"host":             host,
		"packets_sent":     count,
		"packets_received": successful,
		"packet_loss":      fmt.Sprintf("%.1f%%", loss),
		"min_ms":           fmt.Sprintf("%.2f", min),
		"avg_ms":           fmt.Sprintf("%.2f", avg),
		"max_ms":           fmt.Sprintf("%.2f", max),
		"method":           "TCP",
	}

	return result
}

func tryICMPPing(host string, count int) (map[string]interface{}, bool) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", strconv.Itoa(count), host)
	} else {
		cmd = exec.Command("ping", "-c", strconv.Itoa(count), host)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, false
	}

	// Parse básico da saída do ping
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	return map[string]interface{}{
		"host":   host,
		"method": "ICMP",
		"output": strings.TrimSpace(outputStr),
		"lines":  len(lines),
	}, true
}

func testDownloadSpeed() DiagResult {
	result := DiagResult{
		Test:      "download_speed",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	client := &http.Client{Timeout: 60 * time.Second}

	// Usar arquivo de 10MB para teste mais preciso
	testURL := "http://ipv4.download.thinkbroadband.com/10MB.zip"

	start := time.Now()
	resp, err := client.Get(testURL)
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha ao conectar: %v", err)
		return result
	}
	defer resp.Body.Close()

	// Ler todo o conteúdo
	written, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha ao baixar: %v", err)
		return result
	}

	duration := time.Since(start).Seconds()
	// Calcular velocidade em Mbps
	megabytes := float64(written) / (1024 * 1024)
	megabits := megabytes * 8
	speed := megabits / duration

	result.Status = "success"
	result.Details = map[string]interface{}{
		"speed_mbps": fmt.Sprintf("%.2f", speed),
		"size_mb":    fmt.Sprintf("%.2f", megabytes),
		"time_s":     fmt.Sprintf("%.2f", duration),
	}

	return result
}

func testUploadSpeed() DiagResult {
	result := DiagResult{
		Test:      "upload_speed",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Criar dados de teste (1MB)
	testData := make([]byte, 1024*1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	client := &http.Client{Timeout: 60 * time.Second}

	// Usar httpbin para teste de upload
	testURL := "https://httpbin.org/post"

	start := time.Now()
	resp, err := client.Post(testURL, "application/octet-stream", strings.NewReader(string(testData)))
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha ao enviar dados: %v", err)
		return result
	}
	defer resp.Body.Close()

	// Ler resposta
	io.Copy(io.Discard, resp.Body)

	duration := time.Since(start).Seconds()
	megabytes := float64(len(testData)) / (1024 * 1024)
	megabits := megabytes * 8
	speed := megabits / duration

	result.Status = "success"
	result.Details = map[string]interface{}{
		"speed_mbps": fmt.Sprintf("%.2f", speed),
		"size_mb":    fmt.Sprintf("%.2f", megabytes),
		"time_s":     fmt.Sprintf("%.2f", duration),
	}

	return result
}

func testPortConnection(portStr string) DiagResult {
	result := DiagResult{
		Test:      "port",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Verificar se já tem host:porta ou só porta
	target := portStr
	if !strings.Contains(portStr, ":") {
		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			result.Status = "error"
			result.Error = "Porta inválida (deve ser entre 1 e 65535)"
			return result
		}
		target = "localhost:" + portStr
	}

	conn, err := net.DialTimeout("tcp", target, 3*time.Second)
	if err != nil {
		result.Status = "error"
		result.Details = map[string]string{
			"target": target,
			"status": "FECHADA",
		}
		result.Error = err.Error()
		return result
	}
	defer conn.Close()

	result.Status = "success"
	result.Details = map[string]string{
		"target": target,
		"status": "ABERTA",
	}

	return result
}

func testDNSResolution(host string) DiagResult {
	result := DiagResult{
		Test:      "dns",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	start := time.Now()
	ips, err := net.LookupIP(host)
	duration := time.Since(start)

	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha na resolução: %v", err)
		return result
	}

	var ipv4List []string
	var ipv6List []string

	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4List = append(ipv4List, ip.String())
		} else {
			ipv6List = append(ipv6List, ip.String())
		}
	}

	result.Status = "success"
	result.Details = map[string]interface{}{
		"host":        host,
		"ipv4":        ipv4List,
		"ipv6":        ipv6List,
		"response_ms": duration.Milliseconds(),
	}

	return result
}

func testTraceroute(host string) DiagResult {
	result := DiagResult{
		Test:      "traceroute",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("tracert", "-h", "15", host)
	} else {
		cmd = exec.Command("traceroute", "-m", "15", host)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		result.Status = "error"
		result.Error = fmt.Sprintf("Falha ao executar traceroute: %v", err)
		return result
	}

	result.Status = "success"
	result.Details = map[string]interface{}{
		"host":   host,
		"output": strings.TrimSpace(string(output)),
	}

	return result
}

// FUNÇÕES DE SAÍDA
func printResults(results []DiagResult) {
	fmt.Println()
	for _, r := range results {
		printResult(r)
	}
}

func printResult(r DiagResult) {
	switch r.Test {
	case "interfaces":
		if r.Status == "success" {
			if interfaces, ok := r.Details.([]map[string]string); ok {
				printSuccess(fmt.Sprintf("Interfaces de Rede (%d encontradas):", len(interfaces)))
				for _, iface := range interfaces {
					fmt.Printf("  • %s%s%s: %s (Máscara: %s)\n",
						ColorCyan, iface["name"], ColorReset,
						iface["ip"], iface["mask"])
				}
			}
		} else {
			printError("Interfaces: " + r.Error)
		}

	case "public_ip":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]string); ok {
				printSuccess(fmt.Sprintf("IP Público: %s%s%s", ColorBold, details["ip"], ColorReset))
			}
		} else {
			printError("IP Público: " + r.Error)
		}

	case "ping":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]interface{}); ok {
				printSuccess(fmt.Sprintf("Ping para %s (%s):", details["host"], details["method"]))
				fmt.Printf("  • Pacotes: %v enviados, %v recebidos (%s perda)\n",
					details["packets_sent"], details["packets_received"], details["packet_loss"])
				fmt.Printf("  • Latência: min=%s ms, avg=%s ms, max=%s ms\n",
					details["min_ms"], details["avg_ms"], details["max_ms"])
			}
		} else {
			printError("Ping: " + r.Error)
		}

	case "download_speed":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]interface{}); ok {
				printSuccess(fmt.Sprintf("Velocidade de Download: %s%s Mbps%s",
					ColorBold+ColorGreen, details["speed_mbps"], ColorReset))
				fmt.Printf("  • Baixados %s MB em %s segundos\n",
					details["size_mb"], details["time_s"])
			}
		} else {
			printError("Download: " + r.Error)
		}

	case "upload_speed":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]interface{}); ok {
				printSuccess(fmt.Sprintf("Velocidade de Upload: %s%s Mbps%s",
					ColorBold+ColorGreen, details["speed_mbps"], ColorReset))
				fmt.Printf("  • Enviados %s MB em %s segundos\n",
					details["size_mb"], details["time_s"])
			}
		} else {
			printError("Upload: " + r.Error)
		}

	case "port":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]string); ok {
				printSuccess(fmt.Sprintf("Porta %s: %s%s%s",
					details["target"], ColorGreen, details["status"], ColorReset))
			}
		} else {
			if details, ok := r.Details.(map[string]string); ok {
				printError(fmt.Sprintf("Porta %s: %s",
					details["target"], details["status"]))
			} else {
				printError("Porta: " + r.Error)
			}
		}

	case "dns":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]interface{}); ok {
				printSuccess(fmt.Sprintf("Resolução DNS para %s (tempo: %v ms):",
					details["host"], details["response_ms"]))
				if ipv4, ok := details["ipv4"].([]string); ok && len(ipv4) > 0 {
					fmt.Printf("  • IPv4: %s\n", strings.Join(ipv4, ", "))
				}
				if ipv6, ok := details["ipv6"].([]string); ok && len(ipv6) > 0 {
					fmt.Printf("  • IPv6: %s\n", strings.Join(ipv6, ", "))
				}
			}
		} else {
			printError("DNS: " + r.Error)
		}

	case "traceroute":
		if r.Status == "success" {
			if details, ok := r.Details.(map[string]interface{}); ok {
				printSuccess(fmt.Sprintf("Traceroute para %s:", details["host"]))
				fmt.Println(details["output"])
			}
		} else {
			printError("Traceroute: " + r.Error)
		}
	}

	fmt.Println()
}

func outputJSON(results []DiagResult) {
	data := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   VERSION,
		"results":   results,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(`{"error": "Falha ao gerar JSON"}`)
		return
	}
	fmt.Println(string(jsonData))
}
