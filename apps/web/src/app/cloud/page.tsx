'use client';

import { useState, useEffect } from 'react';

interface Account {
  id: string;
  email: string;
  name: string;
  plan: string;
  status: string;
  api_key: string;
  container_id: string;
}

interface Container {
  id: string;
  name: string;
  status: string;
  port: number;
  memory_limit: string;
  cpu_limit: string;
}

export default function CloudDashboard() {
  const [account, setAccount] = useState<Account | null>(null);
  const [container, setContainer] = useState<Container | null>(null);
  const [backups, setBackups] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:7778';

  useEffect(() => {
    fetchAccount();
  }, []);

  const fetchAccount = async () => {
    try {
      const token = localStorage.getItem('hn_token');
      if (!token) {
        window.location.href = '/cloud/login';
        return;
      }

      const response = await fetch(`${API_URL}/api/cloud/auth/me`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (!response.ok) {
        throw new Error('Failed to fetch account');
      }

      const data = await response.json();
      setAccount(data.account);

      // Fetch container info
      const containerResponse = await fetch(`${API_URL}/api/cloud/container`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (containerResponse.ok) {
        const containerData = await containerResponse.json();
        setContainer(containerData.container);
      }

      // Fetch backups
      const backupResponse = await fetch(`${API_URL}/api/cloud/backup/list`, {
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (backupResponse.ok) {
        const backupData = await backupResponse.json();
        setBackups(backupData.backups || []);
      }

      setLoading(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
      setLoading(false);
    }
  };

  const handleStartContainer = async () => {
    try {
      const token = localStorage.getItem('hn_token');
      const response = await fetch(`${API_URL}/api/cloud/container/start`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (response.ok) {
        fetchAccount();
      }
    } catch (err) {
      setError('Failed to start container');
    }
  };

  const handleStopContainer = async () => {
    try {
      const token = localStorage.getItem('hn_token');
      const response = await fetch(`${API_URL}/api/cloud/container/stop`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (response.ok) {
        fetchAccount();
      }
    } catch (err) {
      setError('Failed to stop container');
    }
  };

  const handleCreateBackup = async () => {
    try {
      const token = localStorage.getItem('hn_token');
      const response = await fetch(`${API_URL}/api/cloud/backup`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (response.ok) {
        fetchAccount();
      }
    } catch (err) {
      setError('Failed to create backup');
    }
  };

  const handleCopyApiKey = () => {
    if (account?.api_key) {
      navigator.clipboard.writeText(account.api_key);
      alert('API key copied to clipboard');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-white">Loading...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-900 flex items-center justify-center">
        <div className="text-red-500">{error}</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <header className="bg-gray-800 border-b border-gray-700">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold bg-gradient-to-r from-green-400 to-blue-500 bg-clip-text text-transparent">
            HyperNexus Cloud
          </h1>
          <button
            onClick={() => {
              localStorage.removeItem('hn_token');
              window.location.href = '/cloud/login';
            }}
            className="px-4 py-2 bg-gray-700 rounded hover:bg-gray-600"
          >
            Logout
          </button>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        {/* Account Info */}
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">Account Information</h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-gray-400 text-sm">Name</label>
              <p className="text-lg">{account?.name}</p>
            </div>
            <div>
              <label className="text-gray-400 text-sm">Email</label>
              <p className="text-lg">{account?.email}</p>
            </div>
            <div>
              <label className="text-gray-400 text-sm">Plan</label>
              <p className="text-lg capitalize">{account?.plan}</p>
            </div>
            <div>
              <label className="text-gray-400 text-sm">Status</label>
              <p className="text-lg">
                <span className={`inline-block px-2 py-1 rounded text-sm ${
                  account?.status === 'active' ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300'
                }`}>
                  {account?.status}
                </span>
              </p>
            </div>
          </div>

          <div className="mt-4">
            <label className="text-gray-400 text-sm">API Key</label>
            <div className="flex items-center gap-2 mt-1">
              <code className="bg-gray-700 px-3 py-2 rounded text-sm font-mono flex-1 truncate">
                {account?.api_key}
              </code>
              <button
                onClick={handleCopyApiKey}
                className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500"
              >
                Copy
              </button>
            </div>
          </div>
        </div>

        {/* Container Status */}
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">Container Status</h2>
          {container ? (
            <div>
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-gray-400 text-sm">Container ID</label>
                  <p className="text-lg font-mono">{container.id}</p>
                </div>
                <div>
                  <label className="text-gray-400 text-sm">Status</label>
                  <p className="text-lg">
                    <span className={`inline-block px-2 py-1 rounded text-sm ${
                      container.status === 'running' ? 'bg-green-900 text-green-300' : 'bg-yellow-900 text-yellow-300'
                    }`}>
                      {container.status}
                    </span>
                  </p>
                </div>
                <div>
                  <label className="text-gray-400 text-sm">Port</label>
                  <p className="text-lg">{container.port}</p>
                </div>
                <div>
                  <label className="text-gray-400 text-sm">Memory Limit</label>
                  <p className="text-lg">{container.memory_limit}</p>
                </div>
                <div>
                  <label className="text-gray-400 text-sm">CPU Limit</label>
                  <p className="text-lg">{container.cpu_limit}</p>
                </div>
              </div>

              <div className="flex gap-4">
                {container.status === 'running' ? (
                  <button
                    onClick={handleStopContainer}
                    className="px-4 py-2 bg-red-600 rounded hover:bg-red-500"
                  >
                    Stop Container
                  </button>
                ) : (
                  <button
                    onClick={handleStartContainer}
                    className="px-4 py-2 bg-green-600 rounded hover:bg-green-500"
                  >
                    Start Container
                  </button>
                )}
              </div>
            </div>
          ) : (
            <p className="text-gray-400">No container provisioned</p>
          )}
        </div>

        {/* MCP Connection */}
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">MCP Connection</h2>
          <p className="text-gray-400 mb-4">
            Connect to your HyperNexus Cloud instance using the Streamable HTTP transport.
          </p>
          <div className="bg-gray-700 rounded p-4">
            <label className="text-gray-400 text-sm">Endpoint URL</label>
            <code className="block mt-1 text-green-400 font-mono">
              https://cloud.hypernexus.site/mcp/v1
            </code>
          </div>
          <div className="mt-4 bg-gray-700 rounded p-4">
            <label className="text-gray-400 text-sm">Authentication</label>
            <p className="mt-1 text-sm">
              Include your API key in the <code className="bg-gray-600 px-1 rounded">X-API-Key</code> header
              or as a <code className="bg-gray-600 px-1 rounded">Bearer</code> token in the Authorization header.
            </p>
          </div>
        </div>

        {/* Backups */}
        <div className="bg-gray-800 rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Backups</h2>
          <div className="mb-4">
            <button
              onClick={handleCreateBackup}
              className="px-4 py-2 bg-blue-600 rounded hover:bg-blue-500"
            >
              Create Backup
            </button>
          </div>
          {backups.length > 0 ? (
            <div className="space-y-2">
              {backups.map((backup, index) => (
                <div key={index} className="bg-gray-700 rounded p-3 flex justify-between items-center">
                  <span className="font-mono text-sm">{backup}</span>
                  <button className="px-3 py-1 bg-gray-600 rounded text-sm hover:bg-gray-500">
                    Restore
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-gray-400">No backups available</p>
          )}
        </div>
      </main>
    </div>
  );
}
