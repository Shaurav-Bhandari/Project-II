import { createRouter, createWebHistory } from 'vue-router'

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/orders',
        name: 'Orders',
        component: () => import('../views/Orders.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/orders/new',
        name: 'NewOrder',
        component: () => import('../views/NewOrder.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/kot',
        name: 'KOT',
        component: () => import('../views/KOT.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/menu',
        name: 'Menu',
        component: () => import('../views/Menu.vue'),
        meta: { requiresAuth: true, roles: ['admin', 'manager'] }
    },
    {
        path: '/billing',
        name: 'Billing',
        component: () => import('../views/Billing.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/users',
        name: 'Users',
        component: () => import('../views/Users.vue'),
        meta: { requiresAuth: true, roles: ['admin'] }
    },
    {
        path: '/reports',
        name: 'Reports',
        component: () => import('../views/Reports.vue'),
        meta: { requiresAuth: true, roles: ['admin', 'manager'] }
    },
    {
        path: '/tables',
        name: 'Tables',
        component: () => import('../views/Tables.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/staff/order',
        name: 'StaffOrder',
        component: () => import('../views/StaffOrder.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/kitchen',
        name: 'KitchenDisplay',
        component: () => import('../views/KitchenDisplay.vue'),
        meta: { requiresAuth: true, roles: ['kitchen', 'admin', 'manager'], layout: 'none' }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

// Navigation guard
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    const userRole = localStorage.getItem('userRole')

    if (to.meta.requiresAuth && !token) {
        next('/login')
    } else if (to.path === '/login' && token) {
        // After login, kitchen users go straight to KDS
        if (userRole === 'kitchen') {
            next('/kitchen')
        } else {
            next('/')
        }
    } else if (token && userRole === 'kitchen' && to.path !== '/kitchen') {
        // Kitchen users can ONLY access the KDS
        next('/kitchen')
    } else if (to.meta.roles && !to.meta.roles.includes(userRole)) {
        next('/')
    } else {
        next()
    }
})

export default router
