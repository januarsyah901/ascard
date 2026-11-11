// Base URL for the API. Empty = same origin (production, and dev via Vite proxy).
export const API_BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

export const apiUrl = (path) => `${API_BASE}/api${path}`;
