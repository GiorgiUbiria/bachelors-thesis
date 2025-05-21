import { Outlet } from "react-router";
import { useAuth } from "../store/auth";
import NavigationBar from "./NavigationBar";

export default function MainLayout() {
  const { user, logout } = useAuth();
  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-zinc-900 to-purple-950">
      <div className="flex flex-col min-h-screen">
        <NavigationBar user={user} logout={logout} />
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 flex-1">
          <main className="max-w-4xl mx-auto mt-8 p-6 sm:p-8 bg-zinc-900/90 rounded-3xl shadow-2xl purple-glow border border-purple-900">
            <Outlet />
          </main>
        </div>
      </div>
    </div>
  );
} 