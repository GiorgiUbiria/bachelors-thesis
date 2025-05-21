import { useState, useEffect } from "react";
import { useAuth } from "../store/auth";
import { useNavigate } from "react-router";

function Login() {
  const { login, loading, error, user, clearError } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    if (user) navigate("/");
  }, [user, navigate]);

  useEffect(() => {
    clearError();
  }, [clearError]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const success = await login({ email, password });
    if (success) {
      navigate("/");
    }
  };

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-3xl font-bold text-purple-200 mb-8 text-center">Welcome Back</h1>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="space-y-2">
          <label className="block text-sm font-medium text-purple-200">Email</label>
          <input
            type="email"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
            className="w-full px-4 py-2 rounded-lg bg-zinc-800/50 border border-purple-800/50 text-purple-100 placeholder-purple-400/50 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200"
            placeholder="Enter your email"
          />
        </div>
        <div className="space-y-2">
          <label className="block text-sm font-medium text-purple-200">Password</label>
          <input
            type="password"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
            className="w-full px-4 py-2 rounded-lg bg-zinc-800/50 border border-purple-800/50 text-purple-100 placeholder-purple-400/50 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200"
            placeholder="Enter your password"
          />
        </div>
        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 px-4 rounded-lg bg-purple-600 hover:bg-purple-700 text-white font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-zinc-900"
        >
          {loading ? "Logging in..." : "Sign In"}
        </button>
        {error && (
          <div className="mt-4 p-3 rounded-lg bg-red-500/10 border border-red-500/50 text-red-400 text-sm text-center">
            {error}
          </div>
        )}
      </form>
    </div>
  );
}

export default Login; 