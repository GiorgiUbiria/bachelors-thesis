import { useParams } from "react-router";

export default function ProductDetails() {
  const { id } = useParams();
  return <h2>Product Details for ID: {id}</h2>;
} 