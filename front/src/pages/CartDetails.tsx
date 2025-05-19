import { useParams } from "react-router";

export default function CartDetails() {
  const { id } = useParams();
  return <h2>Cart Details for ID: {id}</h2>;
} 