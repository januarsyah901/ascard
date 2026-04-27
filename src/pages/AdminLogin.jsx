import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { Lock, User, AlertCircle } from 'lucide-react';

const AdminLogin = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await fetch('http://localhost:5001/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (response.ok) {
        login(data.token, data.admin);
        navigate('/admin/dashboard');
      } else {
        setError(data.message || 'Login failed');
      }
    } catch (err) {
      setError('Connection error. Is the server running?');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-surface flex items-center justify-center p-4 font-body">
      <div className="max-w-md w-full bg-surface-container-lowest border border-outline p-10 rounded-none shadow-[20px_20px_0px_0px_rgba(0,0,0,0.05)]">
        <div className="mb-10">
          <div className="inline-block bg-primary text-on-primary-fixed px-3 py-1 mb-6 font-label text-[10px] font-bold tracking-[0.2em] uppercase">
            Admin Portal
          </div>
          <h1 className="font-headline font-black text-4xl tracking-tighter uppercase leading-none text-on-surface">
            Clinical <br /> Authority
          </h1>
          <div className="h-1 w-12 bg-primary mt-6"></div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-8">
          {error && (
            <div className="bg-error-container text-on-error-container p-4 rounded-none border-l-4 border-error flex items-center gap-3 font-label text-[11px] uppercase tracking-wider font-bold">
              <AlertCircle className="w-4 h-4 shrink-0" />
              <p>{error}</p>
            </div>
          )}

          <div className="space-y-2">
            <label className="block font-label font-bold text-[10px] uppercase tracking-widest text-on-surface-variant">Access Identity</label>
            <div className="relative">
              <User className="absolute left-0 top-1/2 -translate-y-1/2 w-4 h-4 text-primary" />
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full bg-transparent border-b border-outline text-on-surface pl-8 pr-4 py-3 rounded-none focus:outline-none focus:border-primary transition-all font-body text-sm"
                placeholder="USERNAME"
                required
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="block font-label font-bold text-[10px] uppercase tracking-widest text-on-surface-variant">Security Key</label>
            <div className="relative">
              <Lock className="absolute left-0 top-1/2 -translate-y-1/2 w-4 h-4 text-primary" />
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-transparent border-b border-outline text-on-surface pl-8 pr-4 py-3 rounded-none focus:outline-none focus:border-primary transition-all font-body text-sm"
                placeholder="PASSWORD"
                required
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-primary hover:bg-primary-fixed-dim text-on-primary-fixed font-label font-black text-xs uppercase tracking-[0.3em] py-4 rounded-none transition-all disabled:opacity-50 disabled:cursor-not-allowed mt-4"
          >
            {loading ? 'Verifying...' : 'Authorize Access'}
          </button>
        </form>

        <div className="mt-12 pt-8 border-t border-outline-variant/30">
          <p className="font-label text-[9px] uppercase tracking-[0.2em] text-outline text-center">
            Secured by Clinical Authority v1.0.4
          </p>
        </div>
      </div>
    </div>
  );
};

export default AdminLogin;
