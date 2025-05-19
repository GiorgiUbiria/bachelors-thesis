import { useParams } from "react-router";

export default function UserCart() {
  const { id } = useParams();
  return <h2>User Cart for ID: {id}</h2>;
} 