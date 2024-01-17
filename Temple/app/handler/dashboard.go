package handler

import (
	"fmt"
	slick "temple"

	"templeapp/views/dashboard"
)

func HandleDashboard(c *slick.Context) error {
	fmt.Println(c.Get("requestID"))
	return c.Render(dashboard.Index())
}
