import React from 'react';
import { trackAdClick } from '../../services/api';

const positions = {
  'top-left': 'top-4 left-4',
  'top-right': 'top-4 right-4',
  'bottom-left': 'bottom-16 left-4',
  'bottom-right': 'bottom-16 right-4',
};

const AdOverlay = ({ ad, videoTime }) => {
  const handleClick = async () => {
    window.open(ad.targetUrl, '_blank');

    try {
      await trackAdClick({
        adId: ad.id,
        videoPlaybackTs: videoTime,
      });
    } catch (error) {
      console.error('Failed to track ad click:', error);
    }
  };

  return (
    <div
      className={`absolute ${positions[ad.position]} transition-opacity duration-300`}
      onClick={handleClick}
    >
      <div className="relative group cursor-pointer">
        <img
          src={ad.imageUrl}
          alt="Advertisement"
          className="w-24 h-24 rounded-lg shadow-lg hover:scale-105 transition-transform"
        />
        <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 rounded-lg transition-opacity flex items-center justify-center">
          <span className="text-white text-sm font-medium">Click to Learn More</span>
        </div>
      </div>
    </div>
  );
};

export default AdOverlay;