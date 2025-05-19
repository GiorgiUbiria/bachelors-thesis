import { Outlet } from "react-router";
import { useAuth } from "../store/auth";

export default function MainLayout() {
  const { user, logout } = useAuth();
  return (
    <div>
      <header style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <h1>Bachelor Project</h1>
        {user ? (
          <div>
            <span>Welcome, {user.name || user.email}!</span>
            <button style={{ marginLeft: 8 }} onClick={logout}>Logout</button>
          </div>
        ) : null}
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  );
} 