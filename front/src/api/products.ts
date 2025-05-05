import axios from 'axios';
import { type Product, ProductCategory } from '../types/product';

const API_URL = 'http://localhost:8080/api';

// Get all products
export const getProducts = async (): Promise<Product[]> => {
  try {
    const response = await axios.get(`${API_URL}/products`);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw error;
    }
    throw new Error('Failed to fetch products');
  }
};

// Get a single product by ID
export const getProduct = async (id: number): Promise<Product> => {
  try {
    const response = await axios.get(`${API_URL}/products/${id}`);
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw error;
    }
    throw new Error('Failed to fetch product');
  }
};

// Get products by category
export const getProductsByCategory = async (category: ProductCategory): Promise<Product[]> => {
  const response = await axios.get(`${API_URL}/products/category/${category}`);
  return response.data;
};

// Get popular products
export const getPopularProducts = async (timeRange: '24h' | '7d' | '30d' = '24h'): Promise<{
  time_range: string;
  products: Array<{
    product_id: number;
    views: number;
  }>;
}> => {
  const response = await axios.get(`${API_URL}/analytics/products/popular?timeRange=${timeRange}`);
  return response.data;
};

// Search products
export const searchProducts = async (query: string): Promise<Product[]> => {
  try {
    const response = await axios.get(`${API_URL}/products/search`, {
      params: { query }
    });
    return response.data;
  } catch (error) {
    if (axios.isAxiosError(error)) {
      throw error;
    }
    throw new Error('Failed to search products');
  }
};

// Error handling wrapper
export const handleProductError = (error: any): never => {
  if (error.response) {
    throw new Error(error.response.data.error || 'Failed to fetch products');
  }
  throw new Error('Network error while fetching products');
};

