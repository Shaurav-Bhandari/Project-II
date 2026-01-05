<template>
  <div class="kot-page">
    <!-- Filters -->
    <div class="filters">
      <button 
        v-for="status in statusFilters" 
        :key="status.value"
        :class="['filter-btn', { active: activeFilter === status.value }]"
        @click="activeFilter = status.value"
      >
        {{ status.label }}
        <span v-if="getCount(status.value) > 0" class="filter-count">{{ getCount(status.value) }}</span>
      </button>
    </div>
    
    <!-- KOT Grid -->
    <div class="kot-grid">
      <div v-for="kot in filteredKots" :key="kot.id" class="kot-card">
        <div class="kot-header" :class="kot.status">
          <div>
            <div class="kot-number">KOT #{{ kot.kot_number }}</div>
            <div class="kot-table">Table {{ kot.table_number }} • Order #{{ kot.order_number }}</div>
          </div>
          <span :class="['badge', getStatusBadge(kot.status)]">
            {{ kot.status === 'in_progress' ? 'Cooking' : kot.status }}
          </span>
        </div>
        
        <div class="kot-body">
          <div v-for="item in kot.items" :key="item.id" class="kot-item">
            <div class="kot-item-check">
              <input 
                type="checkbox" 
                :checked="item.status === 'completed'"
                @change="toggleItemStatus(kot.id, item)"
                :id="`item-${item.id}`"
              >
            </div>
            <div class="kot-qty">{{ item.quantity }}x</div>
            <div class="kot-item-details">
              <div :class="['kot-item-name', { completed: item.status === 'completed' }]">
                {{ item.menu_item_name }}
              </div>
              <div v-if="item.special_instructions" class="kot-item-note">
                📝 {{ item.special_instructions }}
              </div>
            </div>
          </div>
        </div>
        
        <div class="kot-meta">
          <div class="kot-time">
            <span class="time-icon">⏱️</span>
            <span>{{ getTimeAgo(kot.created_at) }}</span>
          </div>
          <div v-if="kot.assigned_chef" class="kot-chef">
            <span>👨‍🍳</span>
            <span>{{ kot.assigned_chef }}</span>
          </div>
        </div>
        
        <div class="kot-footer">
          <button 
            v-if="kot.status === 'pending'"
            class="btn btn-warning" 
            style="flex: 1;"
            @click="updateKotStatus(kot, 'in_progress')"
          >
            Start Cooking
          </button>
          <button 
            v-if="kot.status === 'in_progress'"
            class="btn btn-success" 
            style="flex: 1;"
            @click="updateKotStatus(kot, 'completed')"
          >
            Mark Ready
          </button>
          <button 
            v-if="kot.status === 'completed'"
            class="btn btn-secondary" 
            style="flex: 1;"
            disabled
          >
            ✓ Completed
          </button>
        </div>
      </div>
      
      <div v-if="filteredKots.length === 0" class="empty-state" style="grid-column: 1 / -1; padding: 60px;">
        <div class="empty-state-icon">👨‍🍳</div>
        <div class="empty-state-title">No kitchen orders</div>
        <div class="empty-state-text">All caught up! New orders will appear here.</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const activeFilter = ref('all')

const statusFilters = [
  { value: 'all', label: 'All KOTs' },
  { value: 'pending', label: 'Pending' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'completed', label: 'Completed' }
]

const kots = ref([
  {
    id: 1, kot_number: 89, order_number: 1047, table_number: 'T3', status: 'in_progress',
    assigned_chef: 'Chef Marco', created_at: new Date(Date.now() - 12*60000),
    items: [
      { id: 1, menu_item_name: 'Grilled Salmon', quantity: 2, status: 'in_progress', special_instructions: 'Medium rare' },
      { id: 2, menu_item_name: 'Spring Rolls', quantity: 1, status: 'completed', special_instructions: null },
      { id: 3, menu_item_name: 'Fresh Lemonade', quantity: 2, status: 'completed', special_instructions: 'Less ice' }
    ]
  },
  {
    id: 2, kot_number: 90, order_number: 1048, table_number: 'T7', status: 'pending',
    assigned_chef: null, created_at: new Date(Date.now() - 5*60000),
    items: [
      { id: 4, menu_item_name: 'Ribeye Steak', quantity: 1, status: 'pending', special_instructions: 'Well done' },
      { id: 5, menu_item_name: 'Vegetable Pasta', quantity: 2, status: 'pending', special_instructions: null },
      { id: 6, menu_item_name: 'Chocolate Lava Cake', quantity: 2, status: 'pending', special_instructions: null }
    ]
  },
  {
    id: 3, kot_number: 91, order_number: 1049, table_number: 'T1', status: 'pending',
    assigned_chef: null, created_at: new Date(Date.now() - 2*60000),
    items: [
      { id: 7, menu_item_name: 'Chicken Wings', quantity: 2, status: 'pending', special_instructions: 'Extra spicy' },
      { id: 8, menu_item_name: 'Iced Tea', quantity: 2, status: 'pending', special_instructions: null }
    ]
  },
  {
    id: 4, kot_number: 88, order_number: 1046, table_number: 'T5', status: 'completed',
    assigned_chef: 'Chef Anna', created_at: new Date(Date.now() - 25*60000),
    items: [
      { id: 9, menu_item_name: 'Soup of the Day', quantity: 2, status: 'completed', special_instructions: null },
      { id: 10, menu_item_name: 'Cheesecake', quantity: 1, status: 'completed', special_instructions: null }
    ]
  }
])

const filteredKots = computed(() => {
  if (activeFilter.value === 'all') return kots.value
  return kots.value.filter(k => k.status === activeFilter.value)
})

const getCount = (status) => {
  if (status === 'all') return kots.value.length
  return kots.value.filter(k => k.status === status).length
}

const getStatusBadge = (status) => {
  const badges = {
    pending: 'badge-danger',
    in_progress: 'badge-warning',
    completed: 'badge-success'
  }
  return badges[status] || 'badge-primary'
}

const getTimeAgo = (date) => {
  const minutes = Math.floor((Date.now() - new Date(date).getTime()) / 60000)
  if (minutes < 1) return 'Just now'
  if (minutes === 1) return '1 min ago'
  return `${minutes} mins ago`
}

const updateKotStatus = (kot, status) => {
  kot.status = status
  if (status === 'in_progress') {
    kot.assigned_chef = 'Chef Marco'
  }
}

const toggleItemStatus = (kotId, item) => {
  item.status = item.status === 'completed' ? 'in_progress' : 'completed'
}
</script>

<style scoped>
.kot-header.pending { border-left: 4px solid var(--accent-danger); }
.kot-header.in_progress { border-left: 4px solid var(--accent-warning); }
.kot-header.completed { border-left: 4px solid var(--accent-success); }

.filter-count {
  background: var(--accent-danger);
  color: white;
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 10px;
  margin-left: 6px;
}

.kot-item-check {
  margin-right: 8px;
}

.kot-item-check input {
  width: 18px;
  height: 18px;
  accent-color: var(--accent-success);
}

.kot-item-name.completed {
  text-decoration: line-through;
  color: var(--text-muted);
}

.kot-item-details {
  flex: 1;
}

.kot-meta {
  display: flex;
  justify-content: space-between;
  padding: 12px 20px;
  background: var(--bg-secondary);
  font-size: 13px;
  color: var(--text-secondary);
}

.kot-time, .kot-chef {
  display: flex;
  align-items: center;
  gap: 6px;
}
</style>
