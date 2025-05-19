import { useParams } from "react-router";

export default function OrderDetails() {
  const { id } = useParams();
  return <h2>Order Details for ID: {id}</h2>;
} 