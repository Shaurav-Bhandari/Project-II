<template>
  <div class="menu-page">
    <div class="page-actions">
      <div class="filters">
        <button 
          :class="['filter-btn', { active: activeCategory === 'all' }]"
          @click="activeCategory = 'all'"
        >All Items</button>
        <button 
          v-for="cat in categories" 
          :key="cat.id"
          :class="['filter-btn', { active: activeCategory === cat.id }]"
          @click="activeCategory = cat.id"
        >{{ cat.name }}</button>
      </div>
      <div style="display: flex; gap: 12px;">
        <button class="btn btn-secondary" @click="showCategoryModal = true">+ Category</button>
        <button class="btn btn-primary" @click="openItemModal()">+ Menu Item</button>
      </div>
    </div>
    
    <div class="card">
      <div class="table-container">
        <table class="table">
          <thead>
            <tr>
              <th>Item</th>
              <th>Category</th>
              <th>Price</th>
              <th>Prep Time</th>
              <th>Dietary</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in filteredItems" :key="item.id">
              <td>
                <div><strong>{{ item.name }}</strong></div>
                <div style="font-size: 12px; color: var(--text-muted);">{{ item.description }}</div>
              </td>
              <td>{{ getCategoryName(item.category_id) }}</td>
              <td><strong>NRS {{ item.price.toFixed(2) }}</strong></td>
              <td>{{ item.preparation_time }} min</td>
              <td>
                <div class="dietary-tags">
                  <span v-if="item.is_vegetarian" class="badge badge-success">Veg</span>
                  <span v-if="item.is_vegan" class="badge badge-success">Vegan</span>
                  <span v-if="item.is_gluten_free" class="badge badge-info">GF</span>
                  <span v-if="item.spice_level" class="badge badge-warning">🌶️{{ item.spice_level }}</span>
                </div>
              </td>
              <td>
                <span :class="['badge', item.is_available ? 'badge-success' : 'badge-danger']">
                  {{ item.is_available ? 'Available' : 'Unavailable' }}
                </span>
              </td>
              <td>
                <div class="action-btns">
                  <button class="btn btn-sm btn-secondary" @click="openItemModal(item)">Edit</button>
                  <button 
                    :class="['btn', 'btn-sm', item.is_available ? 'btn-warning' : 'btn-success']"
                    @click="toggleAvailability(item)"
                  >
                    {{ item.is_available ? 'Mark Unavailable' : 'Mark Available' }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- Item Modal -->
    <div :class="['modal-overlay', { active: showItemModal }]" @click.self="showItemModal = false">
      <div class="modal" style="max-width: 600px;">
        <div class="modal-header">
          <h3 class="modal-title">{{ editingItem ? 'Edit Menu Item' : 'Add Menu Item' }}</h3>
          <button class="modal-close" @click="showItemModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
            <div class="form-group">
              <label class="form-label">Item Name</label>
              <input v-model="itemForm.name" type="text" class="form-input" required>
            </div>
            <div class="form-group">
              <label class="form-label">Category</label>
              <select v-model="itemForm.category_id" class="form-input form-select">
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Price ($)</label>
              <input v-model.number="itemForm.price" type="number" step="0.01" class="form-input" required>
            </div>
            <div class="form-group">
              <label class="form-label">Prep Time (min)</label>
              <input v-model.number="itemForm.preparation_time" type="number" class="form-input">
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea v-model="itemForm.description" class="form-input" rows="2"></textarea>
          </div>
          <div style="display: flex; gap: 24px; margin-top: 16px;">
            <label style="display: flex; align-items: center; gap: 8px;">
              <input type="checkbox" v-model="itemForm.is_vegetarian"> Vegetarian
            </label>
            <label style="display: flex; align-items: center; gap: 8px;">
              <input type="checkbox" v-model="itemForm.is_vegan"> Vegan
            </label>
            <label style="display: flex; align-items: center; gap: 8px;">
              <input type="checkbox" v-model="itemForm.is_gluten_free"> Gluten-Free
            </label>
          </div>
          <div class="form-group" style="margin-top: 16px;">
            <label class="form-label">Spice Level (0-5)</label>
            <input v-model.number="itemForm.spice_level" type="range" min="0" max="5" class="form-input">
            <div style="text-align: center;">{{ '🌶️'.repeat(itemForm.spice_level) || 'None' }}</div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showItemModal = false">Cancel</button>
          <button class="btn btn-primary" @click="saveItem">{{ editingItem ? 'Update' : 'Add' }} Item</button>
        </div>
      </div>
    </div>
    
    <!-- Category Modal -->
    <div :class="['modal-overlay', { active: showCategoryModal }]" @click.self="showCategoryModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">Add Category</h3>
          <button class="modal-close" @click="showCategoryModal = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Category Name</label>
            <input v-model="categoryForm.name" type="text" class="form-input" required>
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea v-model="categoryForm.description" class="form-input" rows="2"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showCategoryModal = false">Cancel</button>
          <button class="btn btn-primary" @click="saveCategory">Add Category</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { menuAPI } from '../api'

const activeCategory = ref('all')
const showItemModal = ref(false)
const showCategoryModal = ref(false)
const editingItem = ref(null)

const categories = ref([])
const menuItems = ref([])

const itemForm = ref({ name: '', description: '', category_id: '', price: 0, preparation_time: 15, is_vegetarian: false, is_vegan: false, is_gluten_free: false, spice_level: 0 })
const categoryForm = ref({ name: '', description: '' })

const fetchData = async () => {
  try {
    const [catRes, itemRes] = await Promise.all([
      menuAPI.getCategories(),
      menuAPI.getItems()
    ])
    categories.value = catRes.data || []
    menuItems.value = (itemRes.data || []).map(i => ({ ...i, price: Number(i.price) }))
  } catch (e) {
    console.error('Failed to load menu data', e)
  }
}

onMounted(fetchData)

const filteredItems = computed(() => {
  if (activeCategory.value === 'all') return menuItems.value
  return menuItems.value.filter(item => item.category_id === activeCategory.value)
})

const getCategoryName = (id) => categories.value.find(c => c.id === id)?.name || ''

const openItemModal = (item = null) => {
  if (item) {
    editingItem.value = item
    itemForm.value = { ...item }
  } else {
    editingItem.value = null
    itemForm.value = { name: '', description: '', category_id: categories.value[0]?.id || '', price: 0, preparation_time: 15, is_vegetarian: false, is_vegan: false, is_gluten_free: false, spice_level: 0 }
  }
  showItemModal.value = true
}

const saveItem = async () => {
  try {
    if (editingItem.value) {
      await menuAPI.updateItem(editingItem.value.id, itemForm.value)
    } else {
      await menuAPI.createItem(itemForm.value)
    }
    showItemModal.value = false
    await fetchData()
  } catch (e) {
    alert('Failed to save item: ' + (e.response?.data?.error || e.message))
  }
}

const saveCategory = async () => {
  try {
    await menuAPI.createCategory(categoryForm.value)
    categoryForm.value = { name: '', description: '' }
    showCategoryModal.value = false
    await fetchData()
  } catch (e) {
    alert('Failed to save category: ' + (e.response?.data?.error || e.message))
  }
}

const toggleAvailability = async (item) => {
  try {
    await menuAPI.toggleAvailability(item.id)
    item.is_available = !item.is_available
  } catch (e) {
    console.error('Failed to toggle availability', e)
  }
}
</script>

<style scoped>
.dietary-tags { display: flex; gap: 4px; flex-wrap: wrap; }
.action-btns { display: flex; gap: 8px; }
</style>
