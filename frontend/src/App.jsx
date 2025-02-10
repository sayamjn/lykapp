import React, { useState } from 'react';
import { VideoPlayer } from './components';
import AnalyticsDashboard from './components/AnalyticsDashboard/AnalyticsDashboard';

function App() {
  const [showAnalytics, setShowAnalytics] = useState(false);

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">
            Video Player with Ad Overlays
          </h1>
          <button
            onClick={() => setShowAnalytics(!showAnalytics)}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            {showAnalytics ? 'Hide Analytics' : 'Show Analytics'}
          </button>
        </div>

        <VideoPlayer />

        {showAnalytics && (
          <div className="mt-8">
            <AnalyticsDashboard />
          </div>
        )}
      </div>
    </div>
  );
}

export default App;