import { create } from "zustand";
import { LoginInput, RegisterInput, login as apiLogin, register as apiRegister, getUserFromToken, logout as apiLogout } from "../api/users";
import { User } from "../types/user";
import { LoginSchema, RegisterSchema } from "../api/users";

interface AuthState {
  user: User | null;
  token: string | null;
  loading: boolean;
  error: string | null;
  login: (data: LoginInput) => Promise<boolean>;
  register: (data: RegisterInput) => Promise<boolean>;
  logout: () => void;
  setUser: (user: User | null, token: string | null) => void;
  clearError: () => void;
}

export const useAuth = create<AuthState>((set) => ({
  user: getUserFromToken(),
  token: localStorage.getItem("token"),
  loading: false,
  error: null,

  setUser: (user, token) => set({ user, token }),

  clearError: () => set({ error: null }),

  login: async (data) => {
    set({ loading: true, error: null });
    const parsed = LoginSchema.safeParse(data);
    if (!parsed.success) {
      set({ loading: false, error: parsed.error.errors[0].message });
      return false;
    }
    try {
      const res = await apiLogin(data);
      localStorage.setItem("token", res.token);
      set({ user: res.user, token: res.token, loading: false, error: null });
      return true;
    } catch (err: any) {
      const errorMessage = err.error || "Login failed";
      set({ loading: false, error: errorMessage });
      console.error("Login error:", errorMessage);
      return false;
    }
  },

  register: async (data) => {
    set({ loading: true, error: null });
    const parsed = RegisterSchema.safeParse(data);
    if (!parsed.success) {
      set({ loading: false, error: parsed.error.errors[0].message });
      return false;
    }
    try {
      const res = await apiRegister(data);
      localStorage.setItem("token", res.token);
      set({ user: res.user, token: res.token, loading: false, error: null });
      return true;
    } catch (err: any) {
      const errorMessage = err.error || "Registration failed";
      set({ loading: false, error: errorMessage });
      console.error("Registration error:", errorMessage);
      return false;
    }
  },

  logout: () => {
    apiLogout();
    set({ user: null, token: null, error: null });
  },
})); 