package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
	"gopkg.in/resty.v1"
	"log"
	"net/http"
	"os"
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/healthcheck", healthcheck)
	e.GET("/user/info", UserInfo)

	registerServiceWithConsul()
	registerKong()

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func registerKong() {
	fmt.Println("=======START KONG=======")
	client := resty.New()
	res, _ := client.R().
		SetFormData(map[string]string{
			"name": "user-service",
			"path": "/user-service",
			"url": "http://192.168.0.106:3000",
		}).Post("http://localhost:8001/services/")

	fmt.Println(res)
	fmt.Println("=======START KONG=======")
}

// Handler
func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"code":  http.StatusOK,
		"message": "Welcome to User Service",
	})
}

func healthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "Good!")
}

func UserInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"userId":  "1",
		"fullName": "Minh Hoang",
		"avatar": "https://i.pinimg.com/736x/20/b1/0f/20b10fd84beff4fcd2c3dadb8c50ce29.jpg",
		"email": "anhminh10a2hoa@gmail.com",
	})
}

func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "user-service" //replace with service id
	registration.Name = "user-service" //replace with service name
	address := "192.168.0.106"
	registration.Address = address
	if err != nil {
		log.Fatalln(err)
	}

	registration.Port = 3000
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, 3000)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
	fmt.Println("address" + address)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}

func LookupServiceWithConsul(serviceID string) (string, error) {
	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}
	services, err := client.Agent().Services() // array, slice
	if err != nil {
		return "", err
	}

	fmt.Sprint(services)

	srvc := services[serviceID]
	address := "192.168.0.106"
	port := srvc.Port
	return fmt.Sprintf("http://%s:%v", address, port), nil
}