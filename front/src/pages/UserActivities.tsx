import { useParams } from "react-router";

export default function UserActivities() {
  const { id } = useParams();
  return <h2>User Activities for ID: {id}</h2>;
} 