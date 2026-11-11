import React, { createContext, useContext, useState, useEffect } from 'react';
import { apiUrl } from '../lib/api';

const AuthContext = createContext();

const readToken = () => {
  try { return localStorage.getItem('admin_token'); } catch { return null; }
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(readToken);
  const [loading, setLoading] = useState(Boolean(readToken()));

  useEffect(() => {
    if (!token) return;
    let cancelled = false;
    fetch(apiUrl('/auth/me'), { headers: { Authorization: `Bearer ${token}` } })
      .then((res) => (res.ok ? res.json() : Promise.reject(new Error('unauthorized'))))
      .then((data) => { if (!cancelled) setUser(data.admin); })
      .catch(() => {
        if (cancelled) return;
        try { localStorage.removeItem('admin_token'); } catch { /* ignore */ }
        setToken(null);
        setUser(null);
      })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, [token]);

  const login = (newToken, userData) => {
    try { localStorage.setItem('admin_token', newToken); } catch { /* ignore */ }
    setUser(userData);
    setToken(newToken);
  };

  const logout = () => {
    try { localStorage.removeItem('admin_token'); } catch { /* ignore */ }
    setToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => useContext(AuthContext);
