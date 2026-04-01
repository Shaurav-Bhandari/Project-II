<template>
  <div class="staff-order-page">
    <!-- Left: Menu Browser -->
    <div class="menu-browser">
      <!-- Category Tabs -->
      <div class="category-tabs">
        <button
          :class="['cat-tab', { active: selectedCategory === null }]"
          @click="selectedCategory = null"
        >All</button>
        <button
          v-for="cat in categories"
          :key="cat.id"
          :class="['cat-tab', { active: selectedCategory === cat.id }]"
          @click="selectedCategory = cat.id"
        >{{ cat.name }}</button>
      </div>

      <!-- Search -->
      <div class="search-bar">
        <input
          v-model="searchQuery"
          type="text"
          class="form-input"
          placeholder="🔍  Search menu items..."
        >
      </div>

      <!-- Menu Grid -->
      <div class="menu-items-grid">
        <div
          v-for="item in filteredItems"
          :key="item.id"
          :class="['menu-card', { unavailable: !item.is_available }]"
          @click="item.is_available && addToCart(item)"
        >
          <div class="menu-card-badges">
            <span v-if="item.is_vegetarian" class="veg-badge">🌱</span>
            <span v-if="item.is_vegan" class="vegan-badge">🌿</span>
            <span v-if="!item.is_available" class="sold-out-badge">Sold Out</span>
          </div>
          <div class="menu-card-name">{{ item.name }}</div>
          <div class="menu-card-desc">{{ item.description }}</div>
          <div class="menu-card-footer">
            <span class="menu-card-price">NRS {{ item.price.toFixed(0) }}</span>
            <span class="spice-dots" v-if="item.spice_level > 0">
              <span v-for="n in item.spice_level" :key="n" class="spice-dot">🌶️</span>
            </span>
          </div>
          <div class="menu-card-time">⏱️ {{ item.preparation_time }} min</div>
        </div>

        <div v-if="filteredItems.length === 0" class="empty-menu">
          <div style="font-size: 48px; margin-bottom: 12px;">🍽️</div>
          <div>No items found</div>
        </div>
      </div>
    </div>

    <!-- Right: Cart / Order Panel -->
    <div class="cart-panel">
      <div class="cart-header">
        <h3>Current Order</h3>
        <button v-if="cart.length > 0" class="btn btn-sm btn-danger" @click="clearCart">Clear</button>
      </div>

      <!-- Table & Customer -->
      <div class="order-meta">
        <div class="form-group">
          <label class="form-label">Table</label>
          <select v-model="selectedTable" class="form-input form-select">
            <option value="">Takeaway</option>
            <option v-for="table in tables" :key="table.id" :value="table.id">
              {{ table.table_number }} ({{ table.capacity }} seats) — {{ table.status }}
            </option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">Customer Name</label>
          <input v-model="customerName" type="text" class="form-input" placeholder="Optional">
        </div>
      </div>

      <!-- Cart Items -->
      <div class="cart-items" v-if="cart.length > 0">
        <div v-for="item in cart" :key="item.id" class="cart-item">
          <div class="cart-item-info">
            <div class="cart-item-name">{{ item.name }}</div>
            <div class="cart-item-price">NRS {{ item.price.toFixed(0) }} each</div>
          </div>
          <div class="cart-item-controls">
            <button class="qty-btn" @click="updateQty(item, -1)">−</button>
            <span class="qty-display">{{ item.quantity }}</span>
            <button class="qty-btn" @click="updateQty(item, 1)">+</button>
          </div>
          <div class="cart-item-total">NRS {{ (item.price * item.quantity).toFixed(0) }}</div>
        </div>
      </div>
      <div v-else class="cart-empty">
        <div style="font-size: 40px; margin-bottom: 8px;">🛒</div>
        <div>Tap menu items to add</div>
      </div>

      <!-- Notes -->
      <div class="form-group" v-if="cart.length > 0" style="margin-top: 12px;">
        <label class="form-label">Order Notes</label>
        <textarea v-model="orderNotes" class="form-input" rows="2" placeholder="Any special requests..."></textarea>
      </div>

      <!-- Totals -->
      <div class="cart-totals" v-if="cart.length > 0">
        <div class="total-row">
          <span>Subtotal</span>
          <span>NRS {{ subtotal.toFixed(0) }}</span>
        </div>
        <div class="total-row">
          <span>Tax (10%)</span>
          <span>NRS {{ tax.toFixed(0) }}</span>
        </div>
        <div class="total-row total-final">
          <span>Total</span>
          <span>NRS {{ total.toFixed(0) }}</span>
        </div>
      </div>

      <button
        v-if="cart.length > 0"
        class="btn btn-success place-order-btn"
        @click="placeOrder"
        :disabled="isPlacing"
      >
        {{ isPlacing ? 'Placing...' : '🍽️ Place Order' }}
      </button>

      <!-- Recently placed -->
      <div v-if="recentOrder" class="recent-order">
        <div class="recent-order-badge">✅ Order placed & sent to kitchen!</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { menuAPI, tablesAPI, ordersAPI } from '../api'

const categories = ref([])
const menuItems = ref([])
const tables = ref([])
const selectedCategory = ref(null)
const searchQuery = ref('')
const selectedTable = ref('')
const customerName = ref('')
const orderNotes = ref('')
const cart = ref([])
const isPlacing = ref(false)
const recentOrder = ref(null)

// Fetch data on mount
onMounted(async () => {
  try {
    const [catRes, itemRes, tableRes] = await Promise.all([
      menuAPI.getCategories(),
      menuAPI.getItems(),
      tablesAPI.getAll()
    ])
    categories.value = catRes.data || []
    menuItems.value = (itemRes.data || []).map(i => ({ ...i, price: Number(i.price) }))
    tables.value = tableRes.data || []
  } catch (e) {
    console.error('Failed to load data', e)
  }
})

const filteredItems = computed(() => {
  let items = menuItems.value
  if (selectedCategory.value) {
    items = items.filter(i => i.category_id === selectedCategory.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    items = items.filter(i => i.name.toLowerCase().includes(q) || (i.description && i.description.toLowerCase().includes(q)))
  }
  return items
})

const addToCart = (item) => {
  const existing = cart.value.find(c => c.id === item.id)
  if (existing) {
    existing.quantity++
  } else {
    cart.value.push({ ...item, quantity: 1 })
  }
}

const updateQty = (item, delta) => {
  item.quantity += delta
  if (item.quantity <= 0) {
    cart.value = cart.value.filter(c => c.id !== item.id)
  }
}

const clearCart = () => {
  cart.value = []
  orderNotes.value = ''
}

const subtotal = computed(() => cart.value.reduce((s, i) => s + i.price * i.quantity, 0))
const tax = computed(() => subtotal.value * 0.10)
const total = computed(() => subtotal.value + tax.value)

const placeOrder = async () => {
  if (cart.value.length === 0) return
  isPlacing.value = true
  try {
    const orderData = {
      table_id: selectedTable.value || undefined,
      customer_name: customerName.value,
      order_type: selectedTable.value ? 'dine-in' : 'takeaway',
      notes: orderNotes.value,
      items: cart.value.map(i => ({
        menu_item_id: i.id,
        quantity: i.quantity,
        special_instructions: ''
      }))
    }
    const res = await ordersAPI.create(orderData)
    recentOrder.value = res.data
    cart.value = []
    customerName.value = ''
    orderNotes.value = ''
    setTimeout(() => { recentOrder.value = null }, 5000)
  } catch (e) {
    alert('Failed to place order: ' + (e.response?.data?.message || e.message))
  } finally {
    isPlacing.value = false
  }
}
</script>

<style scoped>
.staff-order-page {
  display: grid;
  grid-template-columns: 1fr 380px;
  gap: 24px;
  height: calc(100vh - var(--header-height) - 64px);
}

/* === Category Tabs === */
.category-tabs {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 8px;
  margin-bottom: 16px;
}
.category-tabs::-webkit-scrollbar { height: 4px; }
.category-tabs::-webkit-scrollbar-thumb { background: var(--border-color); border-radius: 2px; }

.cat-tab {
  padding: 10px 20px;
  border-radius: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  white-space: nowrap;
  font-weight: 500;
  font-size: 13px;
  transition: var(--transition-fast);
}
.cat-tab:hover { border-color: var(--accent-primary); color: var(--text-primary); }
.cat-tab.active {
  background: var(--gradient-primary);
  color: white;
  border-color: transparent;
  box-shadow: var(--shadow-glow);
}

/* === Search === */
.search-bar { margin-bottom: 20px; }

/* === Menu Grid === */
.menu-items-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 14px;
  overflow-y: auto;
  max-height: calc(100vh - var(--header-height) - 230px);
  padding-right: 4px;
}
.menu-items-grid::-webkit-scrollbar { width: 4px; }
.menu-items-grid::-webkit-scrollbar-thumb { background: var(--border-color); border-radius: 2px; }

.menu-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  padding: 16px;
  cursor: pointer;
  transition: var(--transition-fast);
  position: relative;
  display: flex;
  flex-direction: column;
}
.menu-card:hover {
  border-color: var(--accent-primary);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}
.menu-card.unavailable {
  opacity: 0.5;
  cursor: not-allowed;
}
.menu-card:active:not(.unavailable) { transform: scale(0.97); }

.menu-card-badges {
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
  min-height: 20px;
}
.veg-badge, .vegan-badge { font-size: 14px; }
.sold-out-badge {
  background: var(--accent-danger);
  color: white;
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 600;
}

.menu-card-name {
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 4px;
}
.menu-card-desc {
  font-size: 12px;
  color: var(--text-muted);
  flex: 1;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.menu-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.menu-card-price {
  font-weight: 700;
  color: var(--accent-primary);
  font-size: 15px;
}
.spice-dots { font-size: 12px; }
.menu-card-time {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 6px;
}

.empty-menu {
  grid-column: 1 / -1;
  text-align: center;
  padding: 60px 20px;
  color: var(--text-muted);
}

/* === Cart Panel === */
.cart-panel {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  padding: 20px;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  max-height: calc(100vh - var(--header-height) - 64px);
}
.cart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.cart-header h3 {
  font-size: 18px;
  font-weight: 700;
}

.order-meta { margin-bottom: 16px; }
.order-meta .form-group { margin-bottom: 10px; }

/* === Cart Items === */
.cart-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  overflow-y: auto;
  max-height: 250px;
  padding-right: 4px;
}
.cart-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: var(--bg-secondary);
  border-radius: var(--border-radius-sm);
}
.cart-item-info { flex: 1; }
.cart-item-name { font-weight: 600; font-size: 13px; }
.cart-item-price { font-size: 11px; color: var(--text-muted); }
.cart-item-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}
.qty-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  color: var(--text-primary);
  font-size: 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-fast);
}
.qty-btn:hover { border-color: var(--accent-primary); color: var(--accent-primary); }
.qty-display {
  font-weight: 700;
  min-width: 20px;
  text-align: center;
}
.cart-item-total {
  font-weight: 700;
  font-size: 13px;
  min-width: 70px;
  text-align: right;
}
.cart-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  padding: 40px 0;
}

/* === Totals === */
.cart-totals {
  border-top: 1px solid var(--border-color);
  margin-top: 16px;
  padding-top: 12px;
}
.total-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 14px;
}
.total-final {
  font-size: 18px;
  font-weight: 700;
  border-top: 1px solid var(--border-color);
  margin-top: 8px;
  padding-top: 12px;
}

.place-order-btn {
  width: 100%;
  margin-top: 16px;
  padding: 16px;
  font-size: 16px;
  font-weight: 700;
}

.recent-order {
  margin-top: 12px;
  text-align: center;
}
.recent-order-badge {
  background: rgba(16, 185, 129, 0.12);
  color: var(--accent-success);
  padding: 10px 16px;
  border-radius: var(--border-radius-sm);
  font-weight: 600;
  font-size: 14px;
  animation: fadeInUp 0.3s ease;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 1024px) {
  .staff-order-page { grid-template-columns: 1fr; height: auto; }
  .menu-items-grid { max-height: 400px; }
  .cart-panel { max-height: none; }
}

@media (max-width: 768px) {
  .menu-items-grid {
    grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
    gap: 10px;
  }
  .menu-card {
    padding: 12px;
  }
  .menu-card-name {
    font-size: 13px;
  }
  .menu-card-price {
    font-size: 14px;
  }
  .cart-header h3 {
    font-size: 16px;
  }
  .cart-totals {
    font-size: 13px;
  }
}
</style>
