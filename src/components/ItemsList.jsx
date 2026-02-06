import { useState, useEffect } from 'react'

function ItemsList({ token, apiUrl, onLogout }) {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchItems()
  }, [])

  const fetchItems = async () => {
    try {
      const response = await fetch(`${apiUrl}/items`)
      const data = await response.json()
      setItems(data || [])
      setLoading(false)
    } catch (error) {
      console.error('Error fetching items:', error)
      setLoading(false)
    }
  }

  const handleAddToCart = async (itemId) => {
    try {
      const response = await fetch(`${apiUrl}/carts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ item_id: itemId }),
      })

      if (response.ok) {
        window.alert('Item added to cart successfully')
      } else {
        window.alert('Failed to add item to cart')
      }
    } catch (error) {
      window.alert('Error adding item to cart')
    }
  }

  const handleCheckout = async () => {
    try {
      const response = await fetch(`${apiUrl}/orders`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (response.ok) {
        window.alert('Order successful')
      } else {
        const error = await response.json()
        window.alert(error.error || 'Failed to checkout')
      }
    } catch (error) {
      window.alert('Error during checkout')
    }
  }

  const handleViewCart = async () => {
    try {
      const response = await fetch(`${apiUrl}/carts`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      const data = await response.json()
      if (data.cart_items && data.cart_items.length > 0) {
        const cartInfo = data.cart_items.map(item =>
          `Cart ID: ${item.cart_id}, Item ID: ${item.item_id}`
        ).join('\n')
        window.alert(cartInfo)
      } else {
        window.alert('Cart is empty')
      }
    } catch (error) {
      window.alert('Error fetching cart')
    }
  }

  const handleViewOrders = async () => {
    try {
      const response = await fetch(`${apiUrl}/orders`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      const data = await response.json()
      if (data && data.length > 0) {
        const orderIds = data.map(order => `Order ID: ${order.id}`).join('\n')
        window.alert(orderIds)
      } else {
        window.alert('No orders found')
      }
    } catch (error) {
      window.alert('Error fetching orders')
    }
  }

  if (loading) {
    return <div className="loading">Loading items...</div>
  }

  return (
    <>
      <header className="header">
        <div className="header-content">
          <h1>Shopping Cart</h1>
          <div className="header-buttons">
            <button className="btn-cart" onClick={handleViewCart}>
              Cart
            </button>
            <button className="btn-checkout" onClick={handleCheckout}>
              Checkout
            </button>
            <button className="btn-orders" onClick={handleViewOrders}>
              Order History
            </button>
          </div>
        </div>
      </header>
      <div className="container">
        <div className="items-grid">
          {items.map((item) => (
            <div
              key={item.id}
              className="item-card"
              onClick={() => handleAddToCart(item.id)}
            >
              <h3>{item.name}</h3>
              <p>{item.description}</p>
              <div className="price">${item.price.toFixed(2)}</div>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}

export default ItemsList
