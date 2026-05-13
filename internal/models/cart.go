package models

type Cart struct {
	ID                 string           `json:"id"`
	Outlet             *Outlet          `json:"outlet"`
	User               *User            `json:"user"`
	Products           []GroceryProduct `json:"products"`
	Total              float64          `json:"total"`
	TotalDiscount      float64
	TotalAfterDiscount float64
}

func (c *Cart) FindProductByID(productID string) (int, *GroceryProduct) {
	for i, p := range c.Products {
		if p.ID == productID {
			return i, &p
		}
	}
	return 0, nil
}

func (c *Cart) CalculateTotal() *Cart {
	var total float64
	for _, p := range c.Products {
		total += float64(p.Qty) * p.MRP
	}
	c.Total = float64(total)
	return c
}

func (c *Cart) CalculateCategoryTotal() map[string]float64 {
	temp := make(map[string]float64)
	for _, p := range c.Products {
		temp[p.Category] += float64(p.Qty) * p.MRP
	}
	return temp
}

func (c *Cart) ApplyDiscount() *Cart {
	data := c.CalculateCategoryTotal()
	totalDiscout := 0.0
	for category, total := range data {
		switch category {
		case "grocery":
			if total > 50.0 {
				totalDiscout += total - (total*5.0)/100
			}
		case "medicine":
			if total > 20.0 {
				totalDiscout += total - (total*10.0)/100
			}
		case "electronic":
			if total > 10.0 {
				totalDiscout += total - (total*15.0)/100
			}
		}
	}

	c.TotalAfterDiscount = totalDiscout
	return c
}
