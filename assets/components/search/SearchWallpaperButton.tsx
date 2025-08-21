import { useComputed, useSignal } from "@preact/signals";
import { ButtonIcon } from "../ButtonIcon.tsx";
import { WallpaperIcon } from "../icons.tsx";
import { jsonFetch } from "../../functions/http.ts";
import clsx from "clsx";

const storageKey = "grafisearch-wallpaper";
const baseWallpaper = "/images/wallpaper.jpg";

type WallpaperState = {
  time: number;
  url: string;
};

// Set the wallpaper from the local storage
const storageValue = localStorage.getItem(storageKey);
const state = storageValue
  ? (JSON.parse(storageValue) as WallpaperState)
  : { time: 0, url: baseWallpaper };
document.body.style.backgroundImage = `url(${state.url})`;

// If the wallpaper is stale, fetch a new one
const today = new Date();
today.setHours(0, 0, 0, 0);
if (today.getDate() !== new Date(state.time).getDate()) {
  applyAnimeWallpaper();
}

// Reset the wallpaper to the base wallpaper
function resetWallpaper() {
  state.url = baseWallpaper;
  state.time = today.getTime();
  localStorage.setItem(storageKey, JSON.stringify(state));
  document.body.style.backgroundImage = `url(${baseWallpaper})`;
}

// Get screen resolution to optimize wallpaper size
function getOptimalMaxSize(): number {
  const screenWidth = window.screen.width;
  const screenHeight = window.screen.height;
  const maxDimension = Math.max(screenWidth, screenHeight);
  
  // Round up to nearest common resolution
  if (maxDimension <= 1920) return 1920;
  if (maxDimension <= 2560) return 2560;
  if (maxDimension <= 3840) return 3840;
  return 4096; // For 4K+ displays
}

async function applyAnimeWallpaper() {
  try {
    const maxSize = getOptimalMaxSize();
    const { wallpaper } = await jsonFetch<{ wallpaper: string }>(
      `/api/wallpaper?max_size=${maxSize}`
    );
    state.url = wallpaper;
    state.time = today.getTime();
    localStorage.setItem(storageKey, JSON.stringify(state));
    document.body.style.backgroundImage = `url(${wallpaper})`;
  } catch (e) {
    console.error("Failed to fetch anime wallpaper:", e);
    // Fallback to base wallpaper on error
    resetWallpaper();
  }
}

/**
 * Toggle button that switch the wallpaper
 */
export function SearchWallpaperButton() {
  const isLoading = useSignal(false);
  const toggleWallpaper = () => {
    if (state.url === baseWallpaper) {
      isLoading.value = true;
      applyAnimeWallpaper().finally(() => (isLoading.value = false));
    } else {
      resetWallpaper();
    }
  };

  const className = useComputed(() =>
    clsx("search-wallpaper", isLoading.value && "is-loading")
  );

  return (
    <ButtonIcon onClick={toggleWallpaper} className={className.value}>
      <WallpaperIcon />
    </ButtonIcon>
  );
}
