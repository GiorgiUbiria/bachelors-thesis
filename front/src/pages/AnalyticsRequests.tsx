import { useEffect, useState } from 'react';
import { fetchRecentRequestLogs, RequestLog } from '../api/analytics';

export default function AnalyticsRequests() {
  const [logs, setLogs] = useState<RequestLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let mounted = true;
    async function load() {
      setLoading(true);
      try {
        const data = await fetchRecentRequestLogs(50);
        if (mounted) setLogs(data);
      } finally {
        if (mounted) setLoading(false);
      }
    }
    load();
    const interval = setInterval(load, 5000);
    return () => {
      mounted = false;
      clearInterval(interval);
    };
  }, []);

  return (
    <div className="p-4">
      <h2 className="text-xl font-bold mb-4">Analytics: Requests</h2>
      {loading ? (
        <div>Loading...</div>
      ) : (
        <table className="min-w-full border text-sm">
          <thead>
            <tr>
              <th className="border px-2">Time</th>
              <th className="border px-2">IP</th>
              <th className="border px-2">Method</th>
              <th className="border px-2">Path</th>
              <th className="border px-2">Status</th>
              <th className="border px-2">Category</th>
              <th className="border px-2">Response Time (ms)</th>
            </tr>
          </thead>
          <tbody>
            {logs.map(log => (
              <tr key={log.id} className={log.category === 'anomaly' ? 'bg-red-100' : log.category === 'warning' ? 'bg-yellow-100' : ''}>
                <td className="border px-2">{new Date(log.created_at).toLocaleString()}</td>
                <td className="border px-2">{log.ip}</td>
                <td className="border px-2">{log.method}</td>
                <td className="border px-2">{log.path}</td>
                <td className="border px-2">{log.status}</td>
                <td className="border px-2 font-bold">{log.category}</td>
                <td className="border px-2">{log.response_time}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
} 