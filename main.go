package main

func main() {
	apiServer := NewAPIServer(":3000")
	apiServer.Run()
}
