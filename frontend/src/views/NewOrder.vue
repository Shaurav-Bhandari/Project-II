<template>
  <div class="new-order-page">
    <div class="order-form-grid">
      <!-- Menu Items Section -->
      <div class="menu-section">
        <!-- Order Info -->
        <div class="card" style="margin-bottom: 24px;">
          <h3 class="card-title" style="margin-bottom: 16px;">Order Details</h3>
          <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px;">
            <div class="form-group" style="margin-bottom: 0;">
              <label class="form-label">Table</label>
              <select v-model="orderInfo.tableId" class="form-input form-select">
                <option value="">Takeaway</option>
                <option v-for="table in availableTables" :key="table.id" :value="table.id">
                  {{ table.table_number }} ({{ table.capacity }} seats)
                </option>
              </select>
            </div>
            <div class="form-group" style="margin-bottom: 0;">
              <label class="form-label">Customer Name</label>
              <input v-model="orderInfo.customerName" type="text" class="form-input" placeholder="Optional">
            </div>
            <div class="form-group" style="margin-bottom: 0;">
              <label class="form-label">Order Type</label>
              <select v-model="orderInfo.orderType" class="form-input form-select">
                <option value="dine-in">Dine In</option>
                <option value="takeaway">Takeaway</option>
                <option value="delivery">Delivery</option>
              </select>
            </div>
          </div>
        </div>
        
        <!-- Category Filter -->
        <div class="filters" style="margin-bottom: 20px;">
          <button 
            :class="['filter-btn', { active: activeCategory === 'all' }]"
            @click="activeCategory = 'all'"
          >
            All Items
          </button>
          <button 
            v-for="cat in categories" 
            :key="cat.id"
            :class="['filter-btn', { active: activeCategory === cat.id }]"
            @click="activeCategory = cat.id"
          >
            {{ cat.name }}
          </button>
        </div>
        
        <!-- Menu Grid -->
        <div class="menu-grid">
          <div 
            v-for="item in filteredMenuItems" 
            :key="item.id"
            :class="['menu-item-card', { unavailable: !item.is_available }]"
            @click="item.is_available && addToCart(item)"
          >
            <div class="menu-item-name">{{ item.name }}</div>
            <div class="menu-item-desc">{{ item.description }}</div>
            <div class="menu-item-price">${{ item.price.toFixed(2) }}</div>
            <div class="menu-item-tags">
              <span v-if="item.is_vegetarian" class="menu-item-tag">🌱 Veg</span>
              <span v-if="item.is_vegan" class="menu-item-tag">🌿 Vegan</span>
              <span v-if="item.is_gluten_free" class="menu-item-tag">🌾 GF</span>
              <span v-if="item.spice_level > 0" class="menu-item-tag">🌶️ {{ '🌶️'.repeat(item.spice_level) }}</span>
            </div>
            <div v-if="!item.is_available" class="unavailable-label">Unavailable</div>
          </div>
        </div>
      </div>
      
      <!-- Cart Section -->
      <div class="order-cart">
        <h3 style="margin-bottom: 20px;">Current Order</h3>
        
        <div v-if="cart.length === 0" class="empty-state" style="padding: 24px 0;">
          <div class="empty-state-icon">🛒</div>
          <div class="empty-state-text">Add items to start order</div>
        </div>
        
        <div v-else>
          <div v-for="item in cart" :key="item.id" class="cart-item">
            <div class="cart-item-qty">
              <button class="qty-btn" @click="decreaseQty(item)">-</button>
              <span>{{ item.quantity }}</span>
              <button class="qty-btn" @click="increaseQty(item)">+</button>
            </div>
            <div class="cart-item-info">
              <div class="cart-item-name">{{ item.name }}</div>
              <div class="cart-item-price">${{ item.price.toFixed(2) }} each</div>
            </div>
            <div class="cart-item-total">${{ (item.price * item.quantity).toFixed(2) }}</div>
          </div>
          
          <div class="cart-summary">
            <div class="cart-row">
              <span>Subtotal</span>
              <span>${{ subtotal.toFixed(2) }}</span>
            </div>
            <div class="cart-row">
              <span>Tax (10%)</span>
              <span>${{ tax.toFixed(2) }}</span>
            </div>
            <div class="cart-row cart-total">
              <span>Total</span>
              <span>${{ total.toFixed(2) }}</span>
            </div>
          </div>
          
          <div class="form-group" style="margin-top: 20px;">
            <label class="form-label">Order Notes</label>
            <textarea v-model="orderInfo.notes" class="form-input" rows="2" placeholder="Special instructions..."></textarea>
          </div>
          
          <div style="display: flex; gap: 12px; margin-top: 20px;">
            <button class="btn btn-secondary" style="flex: 1;" @click="clearCart">Clear</button>
            <button class="btn btn-primary" style="flex: 2;" @click="submitOrder" :disabled="cart.length === 0">
              Place Order
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const orderInfo = ref({
  tableId: '',
  customerName: '',
  orderType: 'dine-in',
  notes: ''
})

const activeCategory = ref('all')
const cart = ref([])

const availableTables = ref([
  { id: 1, table_number: 'T1', capacity: 2 },
  { id: 2, table_number: 'T2', capacity: 2 },
  { id: 4, table_number: 'T4', capacity: 4 },
  { id: 6, table_number: 'T6', capacity: 6 },
  { id: 8, table_number: 'T8', capacity: 4 },
  { id: 10, table_number: 'T10', capacity: 2 }
])

const categories = ref([
  { id: 1, name: 'Appetizers' },
  { id: 2, name: 'Main Courses' },
  { id: 3, name: 'Desserts' },
  { id: 4, name: 'Beverages' }
])

const menuItems = ref([
  { id: 1, category_id: 1, name: 'Spring Rolls', description: 'Crispy vegetable spring rolls', price: 8.99, is_available: true, is_vegetarian: true, is_vegan: false, is_gluten_free: false, spice_level: 1 },
  { id: 2, category_id: 1, name: 'Chicken Wings', description: 'Spicy buffalo wings', price: 12.99, is_available: true, is_vegetarian: false, is_vegan: false, is_gluten_free: true, spice_level: 3 },
  { id: 3, category_id: 1, name: 'Soup of the Day', description: 'Chef\'s special soup', price: 6.99, is_available: true, is_vegetarian: true, is_vegan: true, is_gluten_free: true, spice_level: 0 },
  { id: 4, category_id: 2, name: 'Grilled Salmon', description: 'Atlantic salmon with lemon butter', price: 24.99, is_available: true, is_vegetarian: false, is_vegan: false, is_gluten_free: true, spice_level: 0 },
  { id: 5, category_id: 2, name: 'Ribeye Steak', description: '12oz prime cut', price: 34.99, is_available: true, is_vegetarian: false, is_vegan: false, is_gluten_free: true, spice_level: 0 },
  { id: 6, category_id: 2, name: 'Vegetable Pasta', description: 'Penne with roasted vegetables', price: 16.99, is_available: true, is_vegetarian: true, is_vegan: false, is_gluten_free: false, spice_level: 1 },
  { id: 7, category_id: 2, name: 'Chicken Tikka Masala', description: 'Creamy tomato curry', price: 18.99, is_available: false, is_vegetarian: false, is_vegan: false, is_gluten_free: true, spice_level: 2 },
  { id: 8, category_id: 3, name: 'Chocolate Lava Cake', description: 'Warm cake with molten center', price: 9.99, is_available: true, is_vegetarian: true, is_vegan: false, is_gluten_free: false, spice_level: 0 },
  { id: 9, category_id: 3, name: 'Cheesecake', description: 'NY style with berry compote', price: 8.99, is_available: true, is_vegetarian: true, is_vegan: false, is_gluten_free: false, spice_level: 0 },
  { id: 10, category_id: 4, name: 'Fresh Lemonade', description: 'Homemade fresh', price: 4.99, is_available: true, is_vegetarian: true, is_vegan: true, is_gluten_free: true, spice_level: 0 },
  { id: 11, category_id: 4, name: 'Iced Tea', description: 'Regular or peach', price: 3.99, is_available: true, is_vegetarian: true, is_vegan: true, is_gluten_free: true, spice_level: 0 },
  { id: 12, category_id: 4, name: 'Coffee', description: 'Freshly brewed', price: 3.49, is_available: true, is_vegetarian: true, is_vegan: true, is_gluten_free: true, spice_level: 0 }
])

const filteredMenuItems = computed(() => {
  if (activeCategory.value === 'all') return menuItems.value
  return menuItems.value.filter(item => item.category_id === activeCategory.value)
})

const subtotal = computed(() => cart.value.reduce((sum, item) => sum + (item.price * item.quantity), 0))
const tax = computed(() => subtotal.value * 0.10)
const total = computed(() => subtotal.value + tax.value)

const addToCart = (item) => {
  const existing = cart.value.find(i => i.id === item.id)
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push({ ...item, quantity: 1 })
  }
}

const increaseQty = (item) => item.quantity++
const decreaseQty = (item) => {
  if (item.quantity > 1) {
    item.quantity--
  } else {
    cart.value = cart.value.filter(i => i.id !== item.id)
  }
}

const clearCart = () => {
  cart.value = []
}

const submitOrder = async () => {
  // In production, this would call the API
  alert('Order placed successfully!')
  router.push('/orders')
}
</script>

<style scoped>
.menu-item-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin: 4px 0 8px;
}

.unavailable-label {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0,0,0,0.8);
  color: var(--accent-danger);
  padding: 8px 16px;
  border-radius: var(--border-radius-sm);
  font-weight: 600;
}

.menu-item-card {
  position: relative;
}
</style>
