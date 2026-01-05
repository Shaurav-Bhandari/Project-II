<template>
  <div class="app-container" v-if="isAuthenticated">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <div class="sidebar-logo">
          <div class="logo-icon">🍽️</div>
          <span class="logo-text">RestaurantOS</span>
        </div>
      </div>
      
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div class="nav-section-title">Main</div>
          <router-link to="/" class="nav-item" :class="{ active: $route.path === '/' }">
            <span class="nav-item-icon">📊</span>
            <span>Dashboard</span>
          </router-link>
          <router-link to="/orders" class="nav-item" :class="{ active: $route.path.startsWith('/orders') }">
            <span class="nav-item-icon">📝</span>
            <span>Orders</span>
            <span v-if="activeOrderCount > 0" class="nav-item-badge">{{ activeOrderCount }}</span>
          </router-link>
          <router-link to="/kot" class="nav-item" :class="{ active: $route.path === '/kot' }">
            <span class="nav-item-icon">👨‍🍳</span>
            <span>Kitchen (KOT)</span>
            <span v-if="pendingKotCount > 0" class="nav-item-badge">{{ pendingKotCount }}</span>
          </router-link>
          <router-link to="/tables" class="nav-item" :class="{ active: $route.path === '/tables' }">
            <span class="nav-item-icon">🪑</span>
            <span>Tables</span>
          </router-link>
        </div>
        
        <div class="nav-section">
          <div class="nav-section-title">Management</div>
          <router-link to="/menu" class="nav-item" :class="{ active: $route.path === '/menu' }">
            <span class="nav-item-icon">🍔</span>
            <span>Menu</span>
          </router-link>
          <router-link to="/billing" class="nav-item" :class="{ active: $route.path === '/billing' }">
            <span class="nav-item-icon">💳</span>
            <span>Billing</span>
          </router-link>
        </div>
        
        <div class="nav-section" v-if="isAdmin">
          <div class="nav-section-title">Admin</div>
          <router-link to="/users" class="nav-item" :class="{ active: $route.path === '/users' }">
            <span class="nav-item-icon">👥</span>
            <span>Users</span>
          </router-link>
          <router-link to="/reports" class="nav-item" :class="{ active: $route.path === '/reports' }">
            <span class="nav-item-icon">📈</span>
            <span>Reports</span>
          </router-link>
        </div>
      </nav>
      
      <div class="sidebar-footer" style="padding: 16px; border-top: 1px solid var(--border-color);">
        <button @click="handleLogout" class="btn btn-secondary" style="width: 100%;">
          🚪 Logout
        </button>
      </div>
    </aside>
    
    <!-- Main Content -->
    <main class="main-content">
      <header class="header">
        <div class="header-left">
          <h1 class="page-title">{{ pageTitle }}</h1>
        </div>
        <div class="header-right">
          <button class="header-btn" title="Notifications">
            🔔
            <span v-if="hasNotifications" class="notification-dot"></span>
          </button>
          <button class="header-btn" title="Quick Order" @click="$router.push('/orders/new')">
            ➕
          </button>
          <div class="user-menu">
            <div class="user-avatar">{{ userInitials }}</div>
            <div class="user-info">
              <div class="user-name">{{ userName }}</div>
              <div class="user-role">{{ userRole }}</div>
            </div>
          </div>
        </div>
      </header>
      
      <div class="page-content">
        <router-view />
      </div>
    </main>
  </div>
  
  <!-- Login Page -->
  <router-view v-else />
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { useOrdersStore } from './stores/orders'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const ordersStore = useOrdersStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const userName = computed(() => authStore.userName)
const userInitials = computed(() => authStore.userInitials)
const userRole = computed(() => authStore.userRole)
const isAdmin = computed(() => ['admin', 'manager'].includes(authStore.userRole))

const activeOrderCount = computed(() => ordersStore.activeOrders.length)
const pendingKotCount = computed(() => ordersStore.pendingKots.length)
const hasNotifications = ref(true)

const pageTitle = computed(() => {
  const titles = {
    '/': 'Dashboard',
    '/orders': 'Order Management',
    '/orders/new': 'New Order',
    '/kot': 'Kitchen Orders',
    '/tables': 'Table Management',
    '/menu': 'Menu Management',
    '/billing': 'Billing & Payments',
    '/users': 'User Management',
    '/reports': 'Reports & Analytics'
  }
  return titles[route.path] || 'Dashboard'
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

onMounted(() => {
  if (isAuthenticated.value) {
    ordersStore.fetchOrders()
    ordersStore.fetchKots()
  }
})
</script>
