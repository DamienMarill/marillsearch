import { useState } from "preact/hooks";

type ImageResult = {
  url: string;
  title: string;
  thumbnail: string;
  source: string;
  domain: string;
  width?: number;
  height?: number;
  size?: string;
};

type Props = {
  result: ImageResult;
};

export function ImageItem({ result }: Props) {
  const [imageLoaded, setImageLoaded] = useState(false);
  const [imageError, setImageError] = useState(false);

  const handleImageLoad = () => {
    setImageLoaded(true);
  };

  const handleImageError = () => {
    setImageError(true);
  };

  const openFullImage = () => {
    window.open(result.url, '_blank');
  };

  const visitSource = (e: Event) => {
    e.stopPropagation();
    window.open(result.source, '_blank');
  };

  if (imageError) {
    return null; // Don't render broken images
  }

  return (
    <div class="image-item" onClick={openFullImage}>
      <div class="image-container">
        <img
          src={result.thumbnail || result.url}
          alt={result.title}
          loading="lazy"
          onLoad={handleImageLoad}
          onError={handleImageError}
          class={imageLoaded ? "loaded" : "loading"}
        />
        {!imageLoaded && <div class="image-placeholder">Loading...</div>}
      </div>
      
      <div class="image-info">
        <div class="image-title" title={result.title}>
          {result.title || "Untitled"}
        </div>
        
        <div class="image-meta">
          <span class="image-domain" onClick={visitSource}>
            {result.domain}
          </span>
          {result.width && result.height && (
            <span class="image-dimensions">
              {result.width} × {result.height}
            </span>
          )}
        </div>
      </div>
    </div>
  );
}