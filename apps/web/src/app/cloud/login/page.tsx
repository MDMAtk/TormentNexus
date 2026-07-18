'use client';

import { useState } from 'react';

export default function CloudLogin() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [isRegister, setIsRegister] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:7778';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const endpoint = isRegister ? '/api/cloud/auth/register' : '/api/cloud/auth/login';
      const body = isRegister ? { email, password, name } : { email, password };

      const response = await fetch(`${API_URL}${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Authentication failed');
      }

      // Store token
      localStorage.setItem('hn_token', data.token);

      // Redirect to dashboard
      window.location.href = '/cloud';
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center">
      <div className="bg-gray-800 rounded-lg p-8 w-full max-w-md">
        <h1 className="text-3xl font-bold mb-6 bg-gradient-to-r from-green-400 to-blue-500 bg-clip-text text-transparent text-center">
          HyperNexus Cloud
        </h1>

        <h2 className="text-xl font-semibold mb-4 text-center text-white">
          {isRegister ? 'Create Account' : 'Sign In'}
        </h2>

        {error && (
          <div className="bg-red-900 text-red-300 px-4 py-2 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          {isRegister && (
            <div>
              <label className="block text-gray-400 text-sm mb-1">Name</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-2 bg-gray-700 rounded text-white focus:outline-none focus:ring-2 focus:ring-green-500"
                placeholder="Your name"
                required
              />
            </div>
          )}

          <div>
            <label className="block text-gray-400 text-sm mb-1">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-2 bg-gray-700 rounded text-white focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="you@example.com"
              required
            />
          </div>

          <div>
            <label className="block text-gray-400 text-sm mb-1">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-2 bg-gray-700 rounded text-white focus:outline-none focus:ring-2 focus:ring-green-500"
              placeholder="••••••••"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-gradient-to-r from-green-500 to-blue-500 rounded font-semibold text-white hover:from-green-600 hover:to-blue-600 disabled:opacity-50"
          >
            {loading ? 'Loading...' : (isRegister ? 'Create Account' : 'Sign In')}
          </button>
        </form>

        <div className="mt-6 text-center">
          <button
            onClick={() => setIsRegister(!isRegister)}
            className="text-green-400 hover:text-green-300"
          >
            {isRegister ? 'Already have an account? Sign in' : "Don't have an account? Register"}
          </button>
        </div>

        <div className="mt-8 pt-6 border-t border-gray-700">
          <h3 className="text-sm font-semibold text-gray-400 mb-2">Plans</h3>
          <div className="space-y-2 text-sm text-gray-400">
            <div className="flex justify-between">
              <span>Free</span>
              <span>512MB RAM, 0.5 CPU</span>
            </div>
            <div className="flex justify-between">
              <span>Starter</span>
              <span>1GB RAM, 1 CPU</span>
            </div>
            <div className="flex justify-between">
              <span>Pro</span>
              <span>2GB RAM, 2 CPU</span>
            </div>
            <div className="flex justify-between">
              <span>Enterprise</span>
              <span>4GB RAM, 4 CPU</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
