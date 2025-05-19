import React from "react";
import { useAuth } from "../store/auth";
import { Navigate, useLocation } from "react-router";

export default function RequireAuth({ children }: React.PropsWithChildren) {
  const { user } = useAuth();
  const location = useLocation();
  if (!user) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return <>{children}</>;
} 