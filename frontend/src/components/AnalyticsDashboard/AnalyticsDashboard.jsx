import React, { useEffect, useState } from 'react';
import { fetchClicks } from '../../services/api';

const AnalyticsDashboard = () => {
  const [clicks, setClicks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const loadClicks = async () => {
    try {
      const data = await fetchClicks();
      setClicks(data);
      setError(null);
    } catch (err) {
      setError('Failed to load analytics data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadClicks();
    const interval = setInterval(loadClicks, 5000);
    return () => clearInterval(interval);
  }, [refreshKey]);

  const calculateMetrics = () => {
    const totalClicks = clicks.length;
    const uniqueUsers = new Set(clicks.map(click => click.ipAddress)).size;
    const clicksByAd = clicks.reduce((acc, click) => {
      acc[click.adId] = (acc[click.adId] || 0) + 1;
      return acc;
    }, {});

    const avgClickTime = totalClicks > 0
      ? (clicks.reduce((acc, click) => acc + click.videoPlaybackTs, 0) / totalClicks).toFixed(1)
      : '0.0';

    return { totalClicks, uniqueUsers, clicksByAd, avgClickTime };
  };

  const handleRefresh = () => {
    setLoading(true);
    setRefreshKey(prev => prev + 1);
  };

  if (loading) {
    return (
      <div className="bg-white p-6 rounded-lg shadow-lg">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Analytics Dashboard</h2>
          <button
            onClick={handleRefresh}
            className="px-4 py-2 text-sm bg-blue-100 text-blue-600 rounded hover:bg-blue-200"
            disabled={loading}
          >
            {loading ? 'Refreshing...' : 'Refresh'}
          </button>
        </div>
        <div className="flex items-center justify-center h-48">
          <div className="text-gray-500">Loading analytics...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white p-6 rounded-lg shadow-lg">
        <div className="text-red-500 p-4 rounded-lg bg-red-50">
          {error}
          <button
            onClick={handleRefresh}
            className="ml-4 text-sm text-red-600 underline hover:no-underline"
          >
            Try again
          </button>
        </div>
      </div>
    );
  }

  const { totalClicks, uniqueUsers, clicksByAd, avgClickTime } = calculateMetrics();

  return (
    <div className="bg-white p-6 rounded-lg shadow-lg">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">Analytics Dashboard</h2>
        <button
          onClick={handleRefresh}
          className="px-4 py-2 text-sm bg-blue-100 text-blue-600 rounded hover:bg-blue-200"
          disabled={loading}
        >
          {loading ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <div className="bg-blue-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold">Total Clicks</h3>
          <p className="text-3xl font-bold text-blue-600">{totalClicks}</p>
        </div>
        
        <div className="bg-green-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold">Unique Users</h3>
          <p className="text-3xl font-bold text-green-600">{uniqueUsers}</p>
        </div>
        
        <div className="bg-purple-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold">Avg. Click Time</h3>
          <p className="text-3xl font-bold text-purple-600">
            {avgClickTime}s
          </p>
        </div>
      </div>

      <div className="mt-6">
        <h3 className="text-xl font-semibold mb-4">Clicks by Ad</h3>
        <div className="space-y-2">
          {Object.entries(clicksByAd).map(([adId, count]) => (
            <div key={adId} className="flex items-center">
              <div className="w-full bg-gray-200 rounded-full h-2.5">
                <div 
                  className="bg-blue-600 h-2.5 rounded-full transition-all duration-500" 
                  style={{ width: `${(count / totalClicks) * 100}%` }}
                ></div>
              </div>
              <span className="ml-2 text-sm whitespace-nowrap">
                Ad {adId}: {count} clicks
              </span>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-8">
        <h3 className="text-xl font-semibold mb-4">Recent Activity</h3>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Time
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Ad ID
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Video Time
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  User IP
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {clicks.slice(0, 10).map((click) => (
                <tr key={click.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(click.timestamp).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {click.adId}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {click.videoPlaybackTs.toFixed(1)}s
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {click.ipAddress}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default AnalyticsDashboard;