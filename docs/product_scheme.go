package docs

import models "github.com/elevenia/product-simple/model"

// swagger:response productResponse
type productResponseWrapper struct {
	// in: body
	Body models.Product
}
