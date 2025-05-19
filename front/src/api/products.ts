import { type Product, ProductCategory } from '../types/product';
import { API_ENDPOINTS } from '../config/api';
import api from '../utils/axiosSetup';

export const getProducts = async (): Promise<Product[]> => {
  try {
    const response = await api.get(API_ENDPOINTS.PRODUCTS.LIST);
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getProduct = async (id: number): Promise<Product> => {
  try {
    const response = await api.get(API_ENDPOINTS.PRODUCTS.DETAIL(id));
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getProductsByCategory = async (category: ProductCategory): Promise<Product[]> => {
  const response = await api.get(API_ENDPOINTS.PRODUCTS.CATEGORY(category));
  return response.data;
};

export const getPopularProducts = async (timeRange: '24h' | '7d' | '30d' = '24h'): Promise<{
  time_range: string;
  products: Array<{
    product_id: number;
    views: number;
  }>;
}> => {
  const response = await api.get(`${API_ENDPOINTS.ANALYTICS.POPULAR_PRODUCTS}?timeRange=${timeRange}`);
  return response.data;
};

export const searchProducts = async (query: string): Promise<Product[]> => {
  try {
    const response = await api.get(API_ENDPOINTS.PRODUCTS.SEARCH, {
      params: { query }
    });
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const handleProductError = (error: any): never => {
  if (error.response) {
    throw new Error(error.response.data.error || 'Failed to fetch products');
  }
  throw new Error('Network error while fetching products');
};

