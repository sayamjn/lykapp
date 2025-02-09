const API_BASE_URL = 'http://localhost:8080/api';

export const fetchAds = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/ads`);
    if (!response.ok) {
      throw new Error('Failed to fetch ads');
    }
    return response.json();
  } catch (error) {
    console.error('Error fetching ads:', error);
    throw error;
  }
};

export const trackAdClick = async (data) => {
  try {
    const response = await fetch(`${API_BASE_URL}/ads/click`, {
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