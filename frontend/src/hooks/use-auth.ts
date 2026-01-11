"use client";

import { useState, useEffect, useCallback } from "react";
import { User } from "@/lib/types";
import { api } from "@/lib/api";

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const checkAuth = useCallback(async () => {
    // Check for token in URL (from OAuth callback)
    if (typeof window !== "undefined") {
      const urlParams = new URLSearchParams(window.location.search);
      const tokenFromUrl = urlParams.get("token");
      const error = urlParams.get("error");

      if (error) {
        console.error("OAuth error:", error);
        // Clear URL params
        window.history.replaceState({}, "", window.location.pathname);
        setLoading(false);
        return;
      }

      if (tokenFromUrl) {
        localStorage.setItem("token", tokenFromUrl);
        // Clear URL params
        window.history.replaceState({}, "", window.location.pathname);
      }
    }

    const token = localStorage.getItem("token");
    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const user = await api.getCurrentUser();
      setUser(user);
    } catch {
      localStorage.removeItem("token");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    checkAuth();
  }, [checkAuth]);

  const login = useCallback((token: string, userData: User) => {
    localStorage.setItem("token", token);
    setUser(userData);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem("token");
    setUser(null);
  }, []);

  return { user, loading, login, logout, checkAuth };
}
