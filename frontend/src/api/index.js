import axios from 'axios'

const api = axios.create({
    baseURL: '/api',
    headers: {
        'Content-Type': 'application/json'
    }
})

// Request interceptor for auth token
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// Response interceptor for error handling
api.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            localStorage.removeItem('userRole')
            window.location.href = '/login'
        }
        return Promise.reject(error)
    }
)

// Auth API
export const authAPI = {
    login: (email, password) => api.post('/auth/login', { email, password }),
    logout: () => api.post('/auth/logout'),
    me: () => api.get('/auth/me')
}

// Users API
export const usersAPI = {
    getAll: () => api.get('/users'),
    getById: (id) => api.get(`/users/${id}`),
    create: (data) => api.post('/users', data),
    update: (id, data) => api.put(`/users/${id}`, data),
    delete: (id) => api.delete(`/users/${id}`),
    getRoles: () => api.get('/roles')
}

// Tables API
export const tablesAPI = {
    getAll: () => api.get('/tables'),
    getById: (id) => api.get(`/tables/${id}`),
    create: (data) => api.post('/tables', data),
    update: (id, data) => api.put(`/tables/${id}`, data),
    updateStatus: (id, status) => api.patch(`/tables/${id}/status`, { status }),
    delete: (id) => api.delete(`/tables/${id}`),
    allocate: (groupSize) => api.post('/tables/allocate', { group_size: groupSize })
}

// Menu API
export const menuAPI = {
    getCategories: () => api.get('/menu/categories'),
    createCategory: (data) => api.post('/menu/categories', data),
    updateCategory: (id, data) => api.put(`/menu/categories/${id}`, data),
    deleteCategory: (id) => api.delete(`/menu/categories/${id}`),

    getItems: (categoryId = null) => api.get('/menu/items', { params: { category_id: categoryId } }),
    getItemById: (id) => api.get(`/menu/items/${id}`),
    createItem: (data) => api.post('/menu/items', data),
    updateItem: (id, data) => api.put(`/menu/items/${id}`, data),
    toggleAvailability: (id) => api.patch(`/menu/items/${id}/availability`),
    deleteItem: (id) => api.delete(`/menu/items/${id}`)
}

// Orders API
export const ordersAPI = {
    getAll: (params) => api.get('/orders', { params }),
    getById: (id) => api.get(`/orders/${id}`),
    create: (data) => api.post('/orders', data),
    update: (id, data) => api.put(`/orders/${id}`, data),
    updateStatus: (id, status) => api.patch(`/orders/${id}/status`, { status }),
    addItem: (orderId, data) => api.post(`/orders/${orderId}/items`, data),
    updateItem: (orderId, itemId, data) => api.put(`/orders/${orderId}/items/${itemId}`, data),
    removeItem: (orderId, itemId) => api.delete(`/orders/${orderId}/items/${itemId}`),
    cancel: (id) => api.delete(`/orders/${id}`)
}

// KOT API
export const kotAPI = {
    getAll: (params) => api.get('/kots', { params }),
    getById: (id) => api.get(`/kots/${id}`),
    updateStatus: (id, status) => api.patch(`/kots/${id}/status`, { status }),
    updateItemStatus: (kotId, itemId, status) => api.patch(`/kots/${kotId}/items/${itemId}/status`, { status }),
    assignChef: (id, chef) => api.patch(`/kots/${id}/assign`, { assigned_chef: chef })
}

// Billing API
export const billingAPI = {
    getPayments: (params) => api.get('/payments', { params }),
    getPaymentById: (id) => api.get(`/payments/${id}`),
    createPayment: (data) => api.post('/payments', data),
    getPaymentMethods: () => api.get('/payment-methods'),
    generateReceipt: (orderId) => api.get(`/orders/${orderId}/receipt`),
    refund: (paymentId) => api.post(`/payments/${paymentId}/refund`)
}

// Reports API
export const reportsAPI = {
    getDashboard: () => api.get('/reports/dashboard'),
    getSalesSummary: (params) => api.get('/reports/sales', { params }),
    getPopularItems: (params) => api.get('/reports/popular-items', { params }),
    getRevenueByDate: (params) => api.get('/reports/revenue', { params }),
    getOrdersByStatus: (params) => api.get('/reports/orders-by-status', { params })
}

export default api
