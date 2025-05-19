import { useParams } from "react-router";

export default function ProductsByCategory() {
  const { category } = useParams();
  return <h2>Products in Category: {category}</h2>;
} 