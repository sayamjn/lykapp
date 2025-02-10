const VITE_API_BASE_URL = 'http://localhost:8080/api';

export const fetchAds = async () => {
  try {
    const response = await fetch(`${VITE_API_BASE_URL}/ads`);
    if (!response.ok) {
      throw new Error('Failed to fetch ads');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching ads:', error);
    throw error;
  }
};

export const trackAdClick = async (data) => {
  try {
    const response = await fetch(`${VITE_API_BASE_URL}/ads/click`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });
    
    if (!response.ok) {
      throw new Error('Failed to track ad click');
    }
  } catch (error) {
    console.error('Error tracking ad click:', error);
    throw error;
  }
};

export const fetchClicks = async () => {
  try {
    const response = await fetch(`${VITE_API_BASE_URL}/ads/clicks`);
    if (!response.ok) {
      throw new Error('Failed to fetch click data');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching clicks:', error);
    throw error;
  }
};

export const fetchAnalytics = async (timeframe = 'day') => {
  try {
    const response = await fetch(`${VITE_API_BASE_URL}/analytics?timeframe=${timeframe}`);
    if (!response.ok) {
      throw new Error('Failed to fetch analytics data');
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching analytics:', error);
    throw error;
  }
};