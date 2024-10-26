package main

func main() {

}
func SwapKeysAndValues(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[v] = k
	}
	return res
}
