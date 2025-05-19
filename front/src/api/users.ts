import { z } from "zod";
import api from "../utils/axiosSetup";
import { API_ENDPOINTS } from "../config/api";

// Zod schemas
export const LoginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

export const RegisterSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
  name: z.string().min(2),
});

export type LoginInput = z.infer<typeof LoginSchema>;
export type RegisterInput = z.infer<typeof RegisterSchema>;

export type User = {
  id: number;
  email: string;
  name: string;
  role: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

export async function login(data: LoginInput): Promise<AuthResponse> {
  const res = await api.post(API_ENDPOINTS.AUTH.LOGIN, data);
  return res.data;
}

export async function register(data: RegisterInput): Promise<AuthResponse> {
  const res = await api.post(API_ENDPOINTS.AUTH.REGISTER, data);
  return res.data;
}

export function getUserFromToken(): User | null {
  const token = localStorage.getItem("token");
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return {
      id: payload.user_id,
      email: payload.email,
      name: payload.name || "",
      role: payload.role,
    };
  } catch {
    return null;
  }
}

export function logout() {
  localStorage.removeItem("token");
}
