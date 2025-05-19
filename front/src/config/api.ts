export const API_BASE_URL = 'http://localhost:8080/api';

export const API_ENDPOINTS = {
  AUTH: {
    LOGIN: `${API_BASE_URL}/auth/login`,
    REGISTER: `${API_BASE_URL}/auth/register`,
  },
  USER: {
    LIST: `${API_BASE_URL}/users`,
    PROFILE: (id: number) => `${API_BASE_URL}/users/${id}`,
    ACTIVITIES: (id: number) => `${API_BASE_URL}/users/${id}/activities`,
    FAVORITES: (id: number) => `${API_BASE_URL}/users/${id}/favorites`,
    CART: (id: number) => `${API_BASE_URL}/users/${id}/cart`,
    ORDERS: (id: number) => `${API_BASE_URL}/users/${id}/orders`,
  },
  PRODUCTS: {
    LIST: `${API_BASE_URL}/products`,
    DETAIL: (id: number) => `${API_BASE_URL}/products/${id}`,
    CATEGORY: (category: string) => `${API_BASE_URL}/products/category/${category}`,
    SEARCH: `${API_BASE_URL}/products/search`,
  },
  CART: {
    DETAIL: (id: number) => `${API_BASE_URL}/cart/${id}`,
    ADD_ITEM: (id: number) => `${API_BASE_URL}/cart/${id}/items`,
    UPDATE_ITEM: (cartId: number, itemId: number) => `${API_BASE_URL}/cart/${cartId}/items/${itemId}`,
    REMOVE_ITEM: (cartId: number, itemId: number) => `${API_BASE_URL}/cart/${cartId}/items/${itemId}`,
  },
  ORDERS: {
    LIST: `${API_BASE_URL}/orders`,
    DETAIL: (id: number) => `${API_BASE_URL}/orders/${id}`,
  },
  ANALYTICS: {
    ACTIVITIES: `${API_BASE_URL}/analytics/activities`,
    REQUESTS: `${API_BASE_URL}/analytics/requests`,
    POPULAR_PRODUCTS: `${API_BASE_URL}/analytics/products/popular`,
    ACTIVE_USERS: `${API_BASE_URL}/analytics/users/active`,
  },
  CSRF: {
    TOKEN: `${API_BASE_URL}/csrf-token`,
  },
}; 