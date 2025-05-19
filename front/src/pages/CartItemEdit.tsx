import { useParams } from "react-router";

export default function CartItemEdit() {
  const { id, itemId } = useParams();
  return <h2>Edit Cart Item {itemId} in Cart {id}</h2>;
} 