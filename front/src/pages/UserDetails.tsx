import { useParams } from "react-router";

export default function UserDetails() {
  const { id } = useParams();
  return <h2>User Details for ID: {id}</h2>;
} 