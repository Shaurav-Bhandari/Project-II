import { defineStore } from 'pinia'
import { authAPI } from '../api'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user') || 'null'),
        token: localStorage.getItem('token') || null,
        isAuthenticated: !!localStorage.getItem('token')
    }),

    getters: {
        userRole: (state) => state.user?.role?.name || 'staff',
        userName: (state) => state.user ? `${state.user.first_name} ${state.user.last_name}` : '',
        userInitials: (state) => state.user ? `${state.user.first_name[0]}${state.user.last_name[0]}` : ''
    },

    actions: {
        async login(email, password) {
            try {
                const response = await authAPI.login(email, password)
                const { token, user } = response.data

                this.token = token
                this.user = user
                this.isAuthenticated = true

                localStorage.setItem('token', token)
                localStorage.setItem('user', JSON.stringify(user))
                localStorage.setItem('userRole', user.role?.name || 'staff')

                return { success: true }
            } catch (error) {
                return {
                    success: false,
                    error: error.response?.data?.message || 'Login failed'
                }
            }
        },

        async logout() {
            try {
                await authAPI.logout()
            } catch (e) {
                // Ignore logout errors
            }

            this.token = null
            this.user = null
            this.isAuthenticated = false

            localStorage.removeItem('token')
            localStorage.removeItem('user')
            localStorage.removeItem('userRole')
        },

        async fetchUser() {
            if (!this.token) return

            try {
                const response = await authAPI.me()
                this.user = response.data
                localStorage.setItem('user', JSON.stringify(response.data))
            } catch (error) {
                this.logout()
            }
        }
    }
})
