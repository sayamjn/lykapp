import React, { useRef, useState, useEffect } from 'react';
import AdOverlay from '../AdOverlay/AdOverlay';
import { useAds } from '../../hooks/useAds';

const VideoPlayer = () => {
  const videoRef = useRef(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const { currentAd, fetchAds } = useAds();

  useEffect(() => {
    fetchAds();
  }, [fetchAds]);

  const handleTimeUpdate = () => {
    if (videoRef.current) {
      setCurrentTime(videoRef.current.currentTime);
    }
  };

  const togglePlay = () => {
    if (videoRef.current) {
      if (isPlaying) {
        videoRef.current.pause();
      } else {
        videoRef.current.play();
      }
      setIsPlaying(!isPlaying);
    }
  };

  return (
    <div className="relative w-full max-w-4xl mx-auto">
      <div className="relative aspect-video bg-black rounded-lg overflow-hidden">
        <video
          ref={videoRef}
          className="w-full h-full object-contain"
          onTimeUpdate={handleTimeUpdate}
          onClick={togglePlay}
        >
          <source src="/sample-video.mp4" type="video/mp4" />
          Your browser does not support the video tag.
        </video>
        
        {currentAd && (
          <AdOverlay
            ad={currentAd}
            videoTime={currentTime}
          />
        )}
        
        <div className="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-black/50 to-transparent">
          <button
            onClick={togglePlay}
            className="text-white px-4 py-2 rounded bg-blue-500 hover:bg-blue-600"
          >
            {isPlaying ? 'Pause' : 'Play'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default VideoPlayer;