import { useParams } from "react-router";

export default function UserOrders() {
  const { id } = useParams();
  return <h2>User Orders for ID: {id}</h2>;
} 