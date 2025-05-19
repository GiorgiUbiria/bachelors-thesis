import { useParams } from "react-router";

export default function UserFavorites() {
  const { id } = useParams();
  return <h2>User Favorites for ID: {id}</h2>;
} 