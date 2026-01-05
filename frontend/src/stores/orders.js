import { defineStore } from 'pinia'
import { ordersAPI, kotAPI } from '../api'

export const useOrdersStore = defineStore('orders', {
    state: () => ({
        orders: [],
        currentOrder: null,
        kots: [],
        loading: false,
        error: null
    }),

    getters: {
        activeOrders: (state) => state.orders.filter(o =>
            !['completed', 'cancelled'].includes(o.status)
        ),
        pendingKots: (state) => state.kots.filter(k => k.status === 'pending'),
        inProgressKots: (state) => state.kots.filter(k => k.status === 'in_progress'),
        completedOrders: (state) => state.orders.filter(o => o.status === 'completed'),
        ordersByStatus: (state) => (status) => state.orders.filter(o => o.status === status)
    },

    actions: {
        async fetchOrders(params = {}) {
            this.loading = true
            try {
                const response = await ordersAPI.getAll(params)
                this.orders = response.data
            } catch (error) {
                this.error = error.response?.data?.message || 'Failed to fetch orders'
            } finally {
                this.loading = false
            }
        },

        async fetchOrderById(id) {
            this.loading = true
            try {
                const response = await ordersAPI.getById(id)
                this.currentOrder = response.data
                return response.data
            } catch (error) {
                this.error = error.response?.data?.message || 'Failed to fetch order'
                return null
            } finally {
                this.loading = false
            }
        },

        async createOrder(orderData) {
            this.loading = true
            try {
                const response = await ordersAPI.create(orderData)
                this.orders.unshift(response.data)
                return { success: true, order: response.data }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to create order'
                }
            } finally {
                this.loading = false
            }
        },

        async updateOrderStatus(id, status) {
            try {
                await ordersAPI.updateStatus(id, status)
                const order = this.orders.find(o => o.id === id)
                if (order) order.status = status
                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to update order status'
                }
            }
        },

        async fetchKots(params = {}) {
            this.loading = true
            try {
                const response = await kotAPI.getAll(params)
                this.kots = response.data
            } catch (error) {
                this.error = error.response?.data?.message || 'Failed to fetch KOTs'
            } finally {
                this.loading = false
            }
        },

        async updateKotStatus(id, status) {
            try {
                await kotAPI.updateStatus(id, status)
                const kot = this.kots.find(k => k.id === id)
                if (kot) kot.status = status
                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to update KOT status'
                }
            }
        },

        async updateKotItemStatus(kotId, itemId, status) {
            try {
                await kotAPI.updateItemStatus(kotId, itemId, status)
                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Failed to update item status'
                }
            }
        }
    }
})
