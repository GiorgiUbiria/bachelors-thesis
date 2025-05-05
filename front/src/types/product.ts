export type Product = {
  id: number;
  name: string;
  description: string;
  price: number;
  stock: number;
  category: string;
  image_url: string;
  created_at: string;
  updated_at: string;
};

export type ProductCategory = 
  | 'Electronics'
  | 'Clothing'
  | 'Books'
  | 'Home'
  | 'Sports'
  | 'Beauty'
  | 'Toys'
  | 'Food';


