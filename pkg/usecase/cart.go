package usecase

import (
	"context"
	"ecommerce/pkg/commonhelp/requests.go"
	"ecommerce/pkg/commonhelp/response"
	"ecommerce/pkg/domain"
	interfaces "ecommerce/pkg/repository/interface"
	services "ecommerce/pkg/usecase/interface"
	"log"

	"github.com/pkg/errors"
)

type CartUsecase struct {
	CartRepo interfaces.CartRepo
}

func NewCartUsecase(cartRepo interfaces.CartRepo) services.CartUsecase {
	return &CartUsecase{
		CartRepo: cartRepo,
	}
}

func (c *CartUsecase) AddCartItem(ctx context.Context, body requests.Cartreq) error {
	// a. product_id validate (find product by product_id)
	product, err := c.CartRepo.FindProduct(ctx, uint(body.ProductId))
	if err != nil {
		log.Printf("[AddCartItem] DB error for product_id=%d: %v", body.ProductId, err)
		return errors.New("invalid product")
	}
	log.Printf("[AddCartItem] Product fetched: %+v", product)
	log.Printf("[AddCartItem] Product qty_in_stock: %d", product.Qty_in_stock)
	if product.Qty_in_stock == 0 {
		log.Printf("[AddCartItem] Product out of stock: product_id=%d", body.ProductId)
		return errors.New("product is currently out of stock")
	}

	// a. find user cart with user_id
	cart, err := c.CartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		log.Printf("[AddCartItem] Failed to find user cart: user_id=%d, err=%v", body.UserID, err)
		return errors.New("failed to find user cart")
	}
	// b. if cart doesn't exist; create new cart with user_id
	if cart.Id == 0 {
		cartId, err := c.CartRepo.SaveCart(ctx, body.UserID)
		if err != nil {
			log.Printf("[AddCartItem] Unable to create cart for user: user_id=%d, err=%v", body.UserID, err)
			return errors.New("unable to create cart for this user")
		}
		cart.Id = cartId
		log.Printf("[AddCartItem] New cart created: cart_id=%d for user_id=%d", cart.Id, body.UserID)
	}

	// a. check if product already exists in cart
	cartitem, err := c.CartRepo.FindCartIDNproductId(ctx, cart.Id, uint(body.ProductId))
	if err != nil {
		log.Printf("[AddCartItem] Failed to check cart items: cart_id=%d, product_id=%d, err=%v", cart.Id, body.ProductId, err)
		return errors.Wrap(err, "failed to check cart items")
	}
	// b. if product already exists in cart
	if cartitem.Id != 0 {
		log.Printf("[AddCartItem] Product already exists in cart: cart_id=%d, product_id=%d", cart.Id, body.ProductId)
		return errors.New("product already exists in cart")
	}

	cartItem := domain.CartItem{
		CartID:    cart.Id,
		ProductId: uint(body.ProductId),
	}

	if err := c.CartRepo.AddCartItem(ctx, cartItem); err != nil {
		log.Printf("[AddCartItem] Failed to add item to cart: %+v, err=%v", cartItem, err)
		return errors.Wrap(err, "failed to add item to cart")
	}

	log.Printf("[AddCartItem] Successfully added product_id=%d to cart_id=%d", body.ProductId, cart.Id)
	return nil
}

func (c *CartUsecase) FindUserCart(ctx context.Context, userID int) (domain.Cart, error) {
	cart, err := c.CartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return domain.Cart{}, errors.Wrap(err, "failed to find user cart")
	}
	return cart, nil
}

func (c *CartUsecase) RemoveFromCart(ctx context.Context, body requests.Cartreq) error {
	// a. product_id validate (find product by product_id)
	product, err := c.CartRepo.FindProduct(ctx, uint(body.ProductId))
	if err != nil {
		return errors.New("invalid product")
	}
	// b. check if product exists
	if product.Id == 0 {
		return errors.New("product is unavailable")
	}

	// a. find user cart with user_id
	cart, err := c.CartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return errors.New("user has no cart")
	}
	// b. if cart doesn't exist
	if cart.Id == 0 {
		return errors.New("cannot remove from cart - cart is empty")
	}

	// a. check if product exists in cart
	cartitem, err := c.CartRepo.FindCartIDNproductId(ctx, cart.Id, uint(body.ProductId))
	if err != nil {
		return errors.Wrap(err, "failed to check cart items")
	}
	if cartitem.Id == 0 {
		return errors.New("product does not exist in your cart")
	}

	if err := c.CartRepo.RemoveCartItem(ctx, cartitem.Id); err != nil {
		return errors.Wrap(err, "failed to remove item from cart")
	}

	return nil
}

func (c *CartUsecase) AddQuantity(ctx context.Context, body requests.Addcount) error {
	product, err := c.CartRepo.FindProduct(ctx, uint(body.ProductId))
	if err != nil {
		return errors.New("invalid product")
	}
	// check if product exists
	if product.Id == 0 {
		return errors.New("product is unavailable")
	}

	if body.Count > uint(product.Qty_in_stock) {
		return errors.New("insufficient product quantity in stock")
	}

	cart, err := c.CartRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return errors.New("user has no cart")
	}

	cartitem, err := c.CartRepo.FindCartIDNproductId(ctx, cart.Id, uint(body.ProductId))
	if err != nil {
		return errors.Wrap(err, "failed to check cart items")
	}
	if cartitem.Id == 0 {
		return errors.New("product does not exist in your cart")
	}

	err = c.CartRepo.AddQuantity(ctx, cartitem.Id, body.Count)
	if err != nil {
		return errors.Wrap(err, "failed to update quantity")
	}

	return nil
}

func (c *CartUsecase) FindCartlistByCartID(ctx context.Context, cartID uint) ([]response.Cartres, error) {
	cartitems, err := c.CartRepo.FindCartlistByCartID(ctx, cartID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cart items")
	}
	return cartitems, nil
}
