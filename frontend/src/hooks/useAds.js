import { useState, useEffect, useCallback } from 'react';
import { fetchAds } from '../services/api';

const AD_ROTATION_INTERVAL = 10000;

export const useAds = () => {
  const [ads, setAds] = useState([]);
  const [currentAd, setCurrentAd] = useState(null);
  const [error, setError] = useState(null);

  const loadAds = useCallback(async () => {
    try {
      const fetchedAds = await fetchAds();
      setAds(fetchedAds);
      if (fetchedAds.length > 0) {
        setCurrentAd(fetchedAds[0]);
      }
      setError(null);
    } catch (err) {
      setError('Failed to fetch ads');
      console.error('Error loading ads:', err);
    }
  }, []);

  useEffect(() => {
    loadAds();
  }, [loadAds]);

  useEffect(() => {
    if (ads.length === 0) return;

    const interval = setInterval(() => {
      setCurrentAd((current) => {
        if (!current) return ads[0];
        const currentIndex = ads.findIndex(ad => ad.id === current.id);
        const nextIndex = (currentIndex + 1) % ads.length;
        return ads[nextIndex];
      });
    }, AD_ROTATION_INTERVAL);

    return () => clearInterval(interval);
  }, [ads]);

  return {
    ads,
    currentAd,
    error,
    fetchAds: loadAds,
  };
};